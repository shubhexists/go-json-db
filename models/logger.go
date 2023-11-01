package models

//Logger
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

//According to ChatGPT : ), Creating Options for any interface is a good idea because it gives you the flexibility to add more options later on without breaking the API.
type Options struct{
	Logger
}
