package config

//go:generate go-mockgen -f github.com/go-nacelle/config -i Config -o config_mock_test.go
//go:generate go-mockgen -f github.com/go-nacelle/config -i Logger -o logger_mock_test.go
//go:generate go-mockgen -f github.com/go-nacelle/config -i Sourcer -o sourcer_mock_test.go
//go:generate go-mockgen -f github.com/go-nacelle/config -i FileSystem -o fs_mock_test.go

import (
	"reflect"
	"testing"

	"github.com/aphistic/sweet"
	"github.com/aphistic/sweet-junit"
	"github.com/fatih/structtag"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&ConfigSuite{})
		s.AddSuite(&DefaultTagSetterSuite{})
		s.AddSuite(&DirectorySourcerSuite{})
		s.AddSuite(&DisplayTagSetterSuite{})
		s.AddSuite(&EnvSourcerSuite{})
		s.AddSuite(&EnvTagPrefixerSuite{})
		s.AddSuite(&FileSourcerSuite{})
		s.AddSuite(&FileTagPrefixerSuite{})
		s.AddSuite(&FileTagSetterSuite{})
		s.AddSuite(&GlobSourcerSuite{})
		s.AddSuite(&JSONSuite{})
		s.AddSuite(&LoggingConfigSuite{})
		s.AddSuite(&MultiSourcerSuite{})
		s.AddSuite(&TestEnvSourcerSuite{})
	})
}

//
//

func ensureEquals(sourcer Sourcer, values []string, expected string) {
	val, flag, err := sourcer.Get(values)
	Expect(err).To(BeNil())
	Expect(flag).To(Equal(FlagFound))
	Expect(val).To(Equal(expected))
}

func ensureMatches(sourcer Sourcer, values []string, expected string) {
	val, flag, err := sourcer.Get(values)
	Expect(err).To(BeNil())
	Expect(flag).To(Equal(FlagFound))
	Expect(val).To(MatchJSON(expected))
}

func ensureMissing(sourcer Sourcer, values []string) {
	_, flag, err := sourcer.Get(values)
	Expect(err).To(BeNil())
	Expect(flag).To(Equal(FlagMissing))
}

func gatherTags(obj interface{}, name string) map[string]string {
	var (
		objValue = reflect.Indirect(reflect.ValueOf(obj))
		objType  = objValue.Type()
	)

	return gatherTagsStruct(objValue, objType, name)
}

func gatherTagsStruct(objValue reflect.Value, objType reflect.Type, name string) map[string]string {
	for i := 0; i < objType.NumField(); i++ {
		var (
			field     = objValue.Field(i)
			fieldType = objType.Field(i)
		)

		if fieldType.Anonymous {
			if tags := gatherTagsStruct(field, fieldType.Type, name); tags != nil {
				return tags
			}
		}

		if fieldType.Name == name {
			if tags, ok := getTags(fieldType); ok {
				return decomposeTags(tags)
			}
		}
	}

	return nil
}

func decomposeTags(tags *structtag.Tags) map[string]string {
	fieldTags := map[string]string{}

	for _, name := range tags.Keys() {
		tag, _ := tags.Get(name)
		fieldTags[name] = tag.Name
	}

	return fieldTags
}
