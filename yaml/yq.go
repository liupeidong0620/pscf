package yamllib

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	errors "github.com/pkg/errors"

	yaml "github.com/liupeidong0620/yaml"
	pscfLog "pscf/log"
)

var trimOutput = true
var writeInplace = true
var writeScript = ""
var outputToJSON = false

//var overwriteFlag = false

//var appendFlag = false

//var verbose = true

var docIndex = "0"

func ReadProperty(args []string) (interface{}, error) {
	var path = ""

	if len(args) < 1 {
		return nil, errors.New("Must provide filename")
	} else if len(args) > 1 {
		path = args[1]
	}

	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return nil, errorParsingDocIndex
	}
	var mappedDocs []interface{}
	var dataBucket interface{}
	//var dataBucket yaml.MapSlice
	var currentIndex = 0
	var errorReadingStream = readStream(args[0], func(decoder *yaml.Decoder) error {
		for {
			errorReading := decoder.Decode(&dataBucket)
			if errorReading == io.EOF {
				pscfLog.Log.Debugf("done %v / %v", currentIndex, docIndexInt)
				if !updateAll && currentIndex <= docIndexInt {
					return fmt.Errorf("Asked to process document index %v but there are only %v document(s)", docIndex, currentIndex)
				}
				return nil
			}
			if errorReading != nil {
				return errorReading
			}
			pscfLog.Log.Debugf("processing %v - requested index %v", currentIndex, docIndexInt)
			if updateAll || currentIndex == docIndexInt {
				pscfLog.Log.Debugf("reading %v in index %v", path, currentIndex)
				//fmt.Println("readPath start ....")
				//fmt.Println("dataBucket: ", dataBucket)
				mappedDoc, errorParsing := readPath(dataBucket, path)
				pscfLog.Log.Debugf("%v", mappedDoc)
				if errorParsing != nil {
					return errors.Wrapf(errorParsing, "Error reading path in document index %v", currentIndex)
				}
				mappedDocs = append(mappedDocs, mappedDoc)
			}
			currentIndex = currentIndex + 1
		}
	})

	if errorReadingStream != nil {
		return nil, errorReadingStream
	}

	if !updateAll {
		dataBucket = mappedDocs[0]
	} else {
		dataBucket = mappedDocs
	}

	return dataBucket, nil
}

/* Conversion readProperty() data */
func ConvReadContext(dataBucket interface{}, display bool, getKeys bool) (interface{}, error) {
	if getKeys {
		dataStr := readMapKeys(dataBucket)
		if display == true {
			if len(dataStr) <= 0 {
				fmt.Println("null")
			} else {
				for _, str := range dataStr {
					fmt.Println(str)
				}
			}
		}
		return dataStr, nil
	} else {
		dataStr, err := toString(dataBucket)
		if err != nil {
			return nil, err
		}
		if display == true {
			fmt.Println(dataStr)
		}
		return dataStr, nil
	}
	return nil, nil
}

func readPath(dataBucket interface{}, path string) (interface{}, error) {
	if path == "" {
		pscfLog.Log.Debug("no path")
		return dataBucket, nil
	}
	var paths = parsePath(path)
	//fmt.Println("readPath: ", paths)
	return recurse(dataBucket, paths[0], paths[1:])
}

func NewProperty(args []string) error {
	updatedData, err := newYaml(args)
	if err != nil {
		return err
	}
	dataStr, err := toString(updatedData)
	if err != nil {
		return err
	}
	fmt.Println(dataStr)
	return nil
}

func newYaml(args []string) (interface{}, error) {
	var writeCommands, writeCommandsError = readWriteCommands(args, 2, "Must provide <path_to_update> <value>")
	if writeCommandsError != nil {
		return nil, writeCommandsError
	}

	var dataBucket interface{}
	var isArray = strings.HasPrefix(writeCommands[0].Key.(string), "[")
	if isArray {
		dataBucket = make([]interface{}, 0)
	} else {
		dataBucket = make(yaml.MapSlice, 0)
	}

	for _, entry := range writeCommands {
		path := entry.Key.(string)
		value := entry.Value
		pscfLog.Log.Debugf("setting %v to %v", path, value)
		var paths = parsePath(path)
		dataBucket = updatedChildValue(dataBucket, paths, value)
	}

	return dataBucket, nil
}

func parseDocumentIndex() (bool, int, error) {
	if docIndex == "*" {
		return true, -1, nil
	}
	docIndexInt64, err := strconv.ParseInt(docIndex, 10, 32)
	if err != nil {
		return false, -1, errors.Wrapf(err, "Document index %v is not a integer or *", docIndex)
	}
	return false, int(docIndexInt64), nil
}

type updateDataFn func(dataBucket interface{}, currentIndex int) (interface{}, error)

func mapYamlDecoder(updateData updateDataFn, encoder *yaml.Encoder) yamlDecoderFn {
	return func(decoder *yaml.Decoder) error {
		var dataBucket interface{}
		var errorReading error
		var errorWriting error
		var errorUpdating error
		var currentIndex = 0

		var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
		if errorParsingDocIndex != nil {
			return errorParsingDocIndex
		}

		for {
			pscfLog.Log.Debugf("Read doc %v", currentIndex)
			errorReading = decoder.Decode(&dataBucket)

			if errorReading == io.EOF {
				if !updateAll && currentIndex <= docIndexInt {
					return fmt.Errorf("Asked to process document index %v but there are only %v document(s)", docIndex, currentIndex)
				}
				return nil
			} else if errorReading != nil {
				return errors.Wrapf(errorReading, "Error reading document at index %v, %v", currentIndex, errorReading)
			}
			dataBucket, errorUpdating = updateData(dataBucket, currentIndex)
			if errorUpdating != nil {
				return errors.Wrapf(errorUpdating, "Error updating document at index %v", currentIndex)
			}

			errorWriting = encoder.Encode(dataBucket)

			if errorWriting != nil {
				return errors.Wrapf(errorWriting, "Error writing document at index %v, %v", currentIndex, errorWriting)
			}
			currentIndex = currentIndex + 1
		}
	}
}

func WriteProperty(args []string) error {
	var writeCommands, writeCommandsError = readWriteCommands(args, 3, "Must provide <filename> <path_to_update> <value>")
	if writeCommandsError != nil {
		return writeCommandsError
	}
	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return errorParsingDocIndex
	}

	var updateData = func(dataBucket interface{}, currentIndex int) (interface{}, error) {
		if updateAll || currentIndex == docIndexInt {
			pscfLog.Log.Debugf("Updating doc %v", currentIndex)
			for _, entry := range writeCommands {
				path := entry.Key.(string)
				value := entry.Value
				pscfLog.Log.Debugf("setting %v to %v", path, value)
				var paths = parsePath(path)
				dataBucket = updatedChildValue(dataBucket, paths, value)
			}
		}
		return dataBucket, nil
	}
	return readAndUpdate(os.Stdout, args[0], updateData)
}

func readAndUpdate(stdOut io.Writer, inputFile string, updateData updateDataFn) error {
	var destination io.Writer
	var destinationName string
	if writeInplace && inputFile != "-" {
		var tempFile, err = ioutil.TempFile("", "temp")
		if err != nil {
			return err
		}
		destinationName = tempFile.Name()
		destination = tempFile
		defer func() {
			safelyCloseFile(tempFile)
			safelyRenameFile(tempFile.Name(), inputFile)
		}()
	} else {
		var writer = bufio.NewWriter(stdOut)
		destination = writer
		destinationName = "Stdout"
		defer safelyFlush(writer)
	}
	var encoder = yaml.NewEncoder(destination)
	pscfLog.Log.Debugf("Writing to %v from %v", destinationName, inputFile)
	return readStream(inputFile, mapYamlDecoder(updateData, encoder))
}

func DeleteProperty(args []string) error {
	if len(args) < 2 {
		return errors.New("Must provide <filename> <path_to_delete>")
	}
	var deletePath = args[1]
	var paths = parsePath(deletePath)
	var updateAll, docIndexInt, errorParsingDocIndex = parseDocumentIndex()
	if errorParsingDocIndex != nil {
		return errorParsingDocIndex
	}

	var updateData = func(dataBucket interface{}, currentIndex int) (interface{}, error) {
		if updateAll || currentIndex == docIndexInt {
			pscfLog.Log.Debugf("Deleting path in doc %v", currentIndex)
			return deleteChildValue(dataBucket, paths), nil
		}
		return dataBucket, nil
	}

	return readAndUpdate(os.Stdout, args[0], updateData)
}

func readWriteCommands(args []string, expectedArgs int, badArgsMessage string) (yaml.MapSlice, error) {
	var writeCommands yaml.MapSlice
	if writeScript != "" {
		if err := readData(writeScript, 0, &writeCommands); err != nil {
			return nil, err
		}
	} else if len(args) < expectedArgs {
		return nil, errors.New(badArgsMessage)
	} else {
		writeCommands = make(yaml.MapSlice, 1)
		writeCommands[0] = yaml.MapItem{Key: args[expectedArgs-2], Value: parseValue(args[expectedArgs-1])}
	}
	return writeCommands, nil
}

func parseValue(argument string) interface{} {
	var value, err interface{}
	var inQuotes = len(argument) > 0 && argument[0] == '"'
	if !inQuotes {
		value, err = strconv.ParseInt(argument, 10, 64)
		if err == nil {
			return value
		}
		value, err = strconv.ParseFloat(argument, 64)
		if err == nil {
			return value
		}
		value, err = strconv.ParseBool(argument)
		if err == nil {
			return value
		}
		if argument == "[]" {
			return make([]interface{}, 0)
		}
		return argument
	}
	return argument[1 : len(argument)-1]
}

func readMapKeys(context interface{}) []string {
	var strData []string
	switch context.(type) {
	case string:
	case []interface{}:
		arrayVal := context.([]interface{})
		for i := 0; i < len(arrayVal); i++ {
			arrayData := readMapKeys(arrayVal[i])
			strData = append(strData, arrayData...)
		}
	case yaml.MapSlice:
		mapVals := context.(yaml.MapSlice)
		for idx := range mapVals {
			entry := mapVals[idx]
			//fmt.Println(entry.Key)
			if valueStr, ok := entry.Key.(string); ok {
				strData = append(strData, valueStr)
			}
		}
	}
	return strData
}

func toString(context interface{}) (string, error) {
	/*
		if outputToJSON {
			return jsonToString(context)
		}*/
	return yamlToString(context)
}

func yamlToString(context interface{}) (string, error) {
	switch context.(type) {
	case string:
		return context.(string), nil
	default:
		return marshalContext(context)
	}
}

func marshalContext(context interface{}) (string, error) {
	out, err := yaml.Marshal(context)

	if err != nil {
		return "", errors.Wrap(err, "error printing yaml")
	}

	outStr := string(out)
	// trim the trailing new line as it's easier for a script to add
	// it in if required than to remove it
	if trimOutput {
		return strings.Trim(outStr, "\n "), nil
	}
	return outStr, nil
}

func safelyRenameFile(from string, to string) {
	if renameError := os.Rename(from, to); renameError != nil {
		pscfLog.Log.Warningf("Error renaming from %v to %v, attemting to copy contents", from, to)
		pscfLog.Log.Warning(renameError.Error())
		// can't do this rename when running in docker to a file targeted in a mounted volume,
		// so gracefully degrade to copying the entire contents.
		if copyError := copyFileContents(from, to); copyError != nil {
			pscfLog.Log.Errorf("Failed copying from %v to %v", from, to)
			pscfLog.Log.Errorf(copyError.Error())
		}
	}
}

// thanks https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer safelyCloseFile(in)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer safelyCloseFile(out)
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func safelyFlush(writer *bufio.Writer) {
	if err := writer.Flush(); err != nil {
		pscfLog.Log.Error("Error flushing writer!")
		pscfLog.Log.Error(err.Error())
	}

}
func safelyCloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		pscfLog.Log.Error("Error closing file!")
		pscfLog.Log.Error(err.Error())
	}
}

type yamlDecoderFn func(*yaml.Decoder) error

func readStream(filename string, yamlDecoder yamlDecoderFn) error {
	if filename == "" {
		return errors.New("Must provide filename")
	}

	var stream io.Reader
	if filename == "-" {
		stream = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer safelyCloseFile(file)
		stream = file
	}
	return yamlDecoder(yaml.NewDecoder(stream))
}

func readData(filename string, indexToRead int, parsedData interface{}) error {
	return readStream(filename, func(decoder *yaml.Decoder) error {
		for currentIndex := 0; currentIndex < indexToRead; currentIndex++ {
			errorSkipping := decoder.Decode(parsedData)
			if errorSkipping != nil {
				return errors.Wrapf(errorSkipping, "Error processing document at index %v, %v", currentIndex, errorSkipping)
			}
		}
		return decoder.Decode(parsedData)
	})
}
