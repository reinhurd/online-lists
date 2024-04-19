package helpers

func GetHomeTemplate() string {
	return `
<html>
<head>
	<title>Online lists service</title>
</head>
<body>
	<h1>Online lists service</h1>
	<form action="/headers" method="get">
		<div><strong>/headers</strong> - get headers from default csv <button type="submit">Visit API</button></div>
	</form>
	<form action="/set_csv" method="get">
		<div>
			<strong>/set_csv</strong> - set default csv
			<input type="text" name="filename" placeholder="Filename"/>
			<button type="submit">Visit API</button>
		</div>
	</form>
	<form action="/list_csv" method="get">
		<div><strong>/list_csv</strong> - list all csv files <button type="submit">Visit API</button></div>
	</form>
	<form action="/add" method="get">
		<div>
			<strong>/add</strong> - add value under header
			<input type="text" name="header" placeholder="Header"/>
			<input type="text" name="value" placeholder="Value"/>
			<button type="submit">Visit API</button>
		</div>
	</form>
	<form action="/ya_file" method="get">
		<div>
			<strong>/ya_file</strong> - download file from Yandex Disk 
			<input type="text" name="filenameSaveAs" placeholder="FilenameSaveAs"/>
			<input type="text" name="path" placeholder="Path"/>
			<button type="submit">Visit API</button>
		</div>
	</form>
	<form action="/ya_list" method="get">
		<div><strong>/ya_list</strong> - list files from Yandex Disk <button type="submit">Visit API</button></div>
	</form>
	<form action="/ya_upload" method="get">
		<div>
			<strong>/ya_upload</strong> - upload file to Yandex Disk
			<input type="text" name="filename" placeholder="Filename"/>
			<button type="submit">Visit API</button>
		</div>
	</form>
</body>
</html>
		
`
}
