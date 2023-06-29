package interface_demo

type People interface {
	Speak(string2 string) string
}

type Stu struct {

}

func (s *Stu) Speak(think string) (talk string) {
	if think == "love" {
		talk = "you are a good boy"
	} else {
		talk = "hi"
	}
	return
}

type People1 interface {
	Show()
}

type Student struct {}

func (s *Student) Show() {}

func live() People1 {
	var stu *Student
	return stu
}