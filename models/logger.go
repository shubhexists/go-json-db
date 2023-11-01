package models


type(
	Logger interface{
		Fatal(string, ...interface{})
		Error(string,...interface{})
		Warn(string, ...interface{})
		Info(string,...interface{})
		Debug(string, ...interface{})
		Trace(string,...interface{})
	}
)


type Options struct{
	Logger
}
