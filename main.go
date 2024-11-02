package main

func main() {

	app := App{}
	app.Initialize(Dbuser, DbPasswrod, DbName)
	app.Run("172.17.197.153:5000")

}
