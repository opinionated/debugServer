var datafunction = function(data) {
	var titles = [];
	$.each( data, function( key, val ) {
		var x = JSON.stringify(val);
		var y = JSON.parse(x);
		$.each( y, function( key, val ) {
			//titles.push("<li id='" + key + "'>" + "<a href = 'localhost/debugServer/article/" + val +"'>" + val + "</li>");
		});
	});

	$( "<ul/>", {
	"class": "articleTitles",
	html: titles.join( "" )
	}).appendTo( "body" );
	console.log(JSON.stringify(data));
}

var failfunc = function(data){
	console.log("bad");
	console.log(data);
}