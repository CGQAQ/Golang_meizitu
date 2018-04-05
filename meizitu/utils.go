package meizitu

import "errors"

type dataType interface{}
type CGIterator struct {
	data []dataType
}

func (itertator *CGIterator) Each(callback func(index int, data dataType)){
	for i, v := range itertator.data{
		callback(i, v)
	}
}

type Pushable interface {
	Push(data ...dataType)
	Pop()    (dataType, error)
	Empty()
	Size()		int
	Iterator()  CGIterator
}

type Stack struct {
	data []dataType
	len  int
}

func (s *Stack) Push(data ...dataType) {
	s.data = append(s.data, data...)
	s.len = len(s.data)
}

func (s *Stack) Pop() (data dataType, err error) {
	defer func() {
		if len(s.data) > 1 {
			s.data = s.data[:len(s.data)-1]
			s.len = len(s.data)
		} else if len(s.data) == 1{
			s.data = make([]dataType, 0)
			s.len = 0
		}
	}()
	length := len(s.data)
	if length > 0{
		return s.data[length-1], nil
	}else {
		return nil, errors.New("empty stack can not pop")
	}
}

func (s *Stack) Empty(){
	s.data = []dataType{}
}

func (s *Stack) Size() int{
	return len(s.data)
}

func (s *Stack) Iterator() CGIterator{
	return CGIterator{s.data}
}






type Queue struct {
	data []dataType
	len int
}

func (q *Queue) Push(data ...dataType) {
	q.data = append(q.data, data...)
	q.len = len(q.data)
}

func (q *Queue) Pop() (data dataType, err error){
	defer func() {
		if len(q.data) > 1 {
			q.data = q.data[1:]
			q.len = len(q.data)
		} else if len(q.data) == 1{
			q.data = make([]dataType, 0)
			q.len = 0
		}
	}()
	length := len(q.data)
	if length > 0{
		return q.data[0], nil
	}else {
		return nil, errors.New("empty stack can not pop")
	}
}

func (q *Queue) Empty(){
	q.data = []dataType{}
}

func (q *Queue) Size() int{
	return len(q.data)
}


func (q *Queue) Iterator() CGIterator{
	return CGIterator{q.data}
}
