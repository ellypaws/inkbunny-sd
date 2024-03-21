package utils

import (
	"bufio"
	"github.com/ellypaws/inkbunny-sd/entities"
	"io"
	"os"
	"strings"
)

type Params map[string]PNGChunk
type PNGChunk map[string]string

type Processor func(string) (Params, error)

// AutoSnep is a Processor that parses yaml like raw txt where each two spaces is a new dict
// It's mostly seen in multi-chunk parameter output from AutoSnep
func AutoSnep(text string) (Params, error) {
	var chunks Params = make(Params)
	scanner := bufio.NewScanner(strings.NewReader(text))

	const parameters = "parameters"
	const Software = "Software"

	var png string
	var key string
	for scanner.Scan() {
		line := scanner.Text()

		indentLevel := len(line) - len(strings.TrimLeft(line, " "))

		switch indentLevel {
		case 0:
			if strings.HasSuffix(line, ":") {
				png = strings.TrimSuffix(line, ":")
				chunks[png] = PNGChunk{}
			}
		case 2: // PNG text chunks:
			if strings.TrimSpace(line) == "PNG text chunks:" {
				chunks[png] = make(PNGChunk)
			}
		case 4: // parameters:
			key = strings.TrimSpace(strings.TrimSuffix(line, ":"))
			//switch {
			//case strings.Contains(line, "parameters:"):
			//	currentKey = parameters
			//case strings.Contains(line, "Software:"):
			//	currentKey = Software
			//}
		case 6:
			if len(chunks[png][key]) > 0 {
				chunks[png][key] += "\n"
			}
			chunks[png][key] += line[6:]
		default:
			// Handle other cases as needed, e.g., deeper indentations or different keys
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return chunks, nil
}

func ParseParams(p Params) map[string]entities.TextToImageRequest {
	var request map[string]entities.TextToImageRequest
	for file, chunk := range p {
		if params, ok := chunk["parameters"]; ok {
			r, err := ParameterHeuristics(params)
			if err != nil {
				continue
			}
			if request == nil {
				request = make(map[string]entities.TextToImageRequest)
			}
			request[file] = r
		}
	}
	return request
}

func FileToRequests(file string, processor Processor) (map[string]entities.TextToImageRequest, error) {
	p, err := fileToParams(file, processor)
	if err != nil {
		return nil, err
	}
	return ParseParams(p), nil
}

// fileToParams reads the file and returns the params using a Processor
func fileToParams(file string, processor Processor) (Params, error) {
	f, err := fileToBytes(file)
	if err != nil {
		return nil, err
	}
	return processor(string(f))
}

// fileToBytes reads the file and returns the content as a byte slice
func fileToBytes(file string) ([]byte, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(f)
}
