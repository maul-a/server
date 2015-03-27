function computersaction (action) {
	a = document.getElementsByTagName('input');
	var mas = new Array();var c = new Array();
	for(i=0;i<a.length - 3;i++) {if (a[i].checked) mas.push(a[i])}
	if (mas.length===0) {alert('Вы ничего не выбрали')} else {
		for(i=0;i<mas.length;i++) {s=mas[i].name; s=s.replace(/id-/,""); c.push(s)}
		var s = c.join();	
		var params = action + '=' + encodeURIComponent(s);
		var xmlhttp = getXmlHttp()
		xmlhttp.open("POST", '/active', true);
		xmlhttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
		xmlhttp.send(params);
		document.location.href = '/active';	
	}

}
document.getElementById('poweroffbutton').onclick = function() {
	computersaction('poweroff');
}

document.getElementById('rebootbutton').onclick = function() {
	computersaction('reboot');
}
document.getElementById('sendfilebutton').onclick = function() {
	a = document.getElementsByTagName('input');
	var mas = new Array();var c = new Array();
	for(i=0;i<a.length - 3;i++) {if (a[i].checked) mas.push(a[i])}
	if (mas.length===0) {alert('Вы ничего не выбрали')} else {
		for(i=0;i<mas.length;i++) {s=mas[i].name; s=s.replace(/id-/,""); c.push(s)}
		var s = c.join();
		document.getElementById('myfiles').click();
	}
document.getElementById('myfiles').onchange = function() {
		document.getElementById('sendfile').value = s;
		document.getElementById('form1').submit();
}	
}

	$('input[type="checkbox"]').checkbox();