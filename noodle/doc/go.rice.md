# go.rice

- https://github.com/GeertJohan/go.rice
- support get content from box as string

````go
// find a rice.Box
templateBox, err := rice.FindBox("example-templates")
if err != nil {
	log.Fatal(err)
}
// get file contents as string
templateString, err := templateBox.String("message.tmpl")
if err != nil {
	log.Fatal(err)
}
````