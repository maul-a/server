function computersaction (action) {
	a = document.getElementsByTagName('input');
	var mas = new Array();var c = new Array();
	for(i=0;i<a.length;i++) {if (a[i].checked) mas.push(a[i])}
	if (mas.length===0) {alert('Вы ничего не выбрали')} else {
		for(i=0;i<mas.length;i++) {s=mas[i].name; s=s.replace(/id-/,""); c.push(s)}
		var s = c.join();	
		var params = action + '=' + encodeURIComponent(s);
		var xmlhttp = getXmlHttp()
		xmlhttp.open("POST", '/computers', true);
		xmlhttp.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
		xmlhttp.send(params);
		document.location.href = '/settings';	
	}

}
document.getElementById('poweronbutton').onclick = function() {
	computersaction('poweron')
}
$('input[type="checkbox"]').checkbox();
