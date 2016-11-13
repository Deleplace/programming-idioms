package pigae

import (
	"fmt"
)

// PipelineDecorator is an artificial structure because pipeline is always only 1 object :(
type PipelineDecorator struct {
	// The actual pipeline
	Data interface{}
	// Some helper data passed as "second pipeline"
	Deco interface{}
}

func decorate(data interface{}, deco interface{}) *PipelineDecorator {
	return &PipelineDecorator{
		Data: data,
		Deco: deco,
	}
}

// Thanks to tux21b
// (see https://stackoverflow.com/questions/18276173/calling-a-template-with-several-pipeline-parameters)
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// Another attempt (2013-12)
// Put the context data in the root pipeline
/*
type CompositeData map[string]interface{}

type LeanPipe struct {
	root   *LeanPipe
	Data   CompositeData
}

func NewLeanPipe(data map[string]interface{}) *LeanPipe {
	lp := &LeanPipe{
		root:   nil,
		Data:   data,
	}
	lp.root = lp
	return lp
}

func (pipe *LeanPipe) Sub(field string) *LeanPipe {
	subdata := pipe.Data[field].(CompositeData)
	return &LeanPipe{
		root:   pipe.root,
		Data:   subdata,
	}
}

func (pipe *LeanPipe) Context(field string) interface{} {
	return pipe.root.Data[field]
}
*/
