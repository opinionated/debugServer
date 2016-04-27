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

var buildArticle = function(data){
	// builds just one article template
	var content = document.querySelector('#article_template').content;

	var titleTag = content.querySelector('#title');
	titleTag.textContent = data["Title"];
	
	var bodyTag = content.querySelector('#body');
	bodyTag.textContent = data["Body"];

	return content;
}

var relatedClicked = function(which){
	console.log(which)
	
	$.getJSON('http://localhost:8002/api/article/' + which, function(data){
		var body = document.querySelector("#body");
		var related_article = body.querySelector("#related_article");

		var bodyTag = related_article.querySelector("#body");
		bodyTag.textContent = data['Body'];

		var titleTag = related_article.querySelector("#title");
		titleTag.textContent = data['Title'];	
	});

}

var getRelated = function(arr){
	for(var index in arr){
		var related = arr[index];
	$.getJSON('http://localhost:8002/api/article/' + related['Title'], 
		function(data){

			// var article = buildArticle(data);
			
			var toAppend = "<button onclick='relatedClicked(\"";
			toAppend += data["Title"];
			toAppend += "\")'>";
			toAppend += data["Title"]
			toAppend += "</button>";
			$("#related_list").append(toAppend)

		})
	}
}

var buildArticlePage = function(data){
	// builds the article page, places divs properly etc

	console.log(data)

	var mainArticle = buildArticle(data);
	var mainDoc = document.querySelector("#body");


	mainDoc.querySelector('#main_article').appendChild(
		document.importNode(mainArticle, true));
	mainDoc.querySelector('#related_article').appendChild(
		document.importNode(mainArticle, true));

	getRelated(data["Related"]);

}

var getArticle = function(after){
	params = getUrlVars()
	which = params["which"]
	// after should always be a callback to buildArticlePage
	$.getJSON('http://localhost:8002/api/article/' + which, after)
}

function getUrlVars() {
	var vars = {};
	var parts = window.location.href.replace(/[?&]+([^=&]+)=([^&]*)/gi, 
		function(m,key,value) {
			vars[key] = value;
		}
	);
	return vars;
}