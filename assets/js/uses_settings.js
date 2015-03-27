document.getElementById('addcomputer').onclick = function() {
if (document.getElementById('addnewcomputer').style.display === "block")
	document.getElementById('addnewcomputer').style.display = "none";
else
	document.getElementById('addnewcomputer').style.display = "block";
}
document.getElementById('deletecomputer').onclick = function() {
	a = document.getElementsByTagName('input'); 
	var s=""; var c = new Array(); var mas=new Array(); 
	var i=0;
	for(i=0;i<a.length - 3;i++) {if (a[i].checked) mas.push(a[i])}
	if (mas.length===0) {alert('Вы ничего не выбрали')} else {
	for(i=0;i<mas.length;i++) {s=mas[i].name; s=s.replace(/id-/,""); c.push(s)}
	var s = c.join();
	var params = 'str=' + encodeURIComponent(s);
	var xmlhttp = getXmlHttp()
	xmlhttp.open("POST", '/delsettings', true);
	xmlhttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
	xmlhttp.send(params);
	document.location.href = '/settings';
	}
}

$('input[type="checkbox"]').checkbox();