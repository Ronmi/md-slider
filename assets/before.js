(function(pres, append, codepen, i, pre, code){
    for (i=0; i<pres.length; i++) {
	pre = pres.item(i)

	if (pre.children.length < 1) continue;

	code = pre.children[0];
	if (code.tagName != "CODE") continue;

	// 把 codeblock 通通加上 prism plugin 的 class
	if (pre.className != "") pre.className += " ";
	pre.className += "line-numbers";

	// 處理各種特殊 code blocks
	switch (code.className) {
	    case "language-html":
		append(pre, codepen("html", code));
		break;
	    case "language-js":
		append(pre, codepen("js", code));
		break;
	    case "language-css":
		append(pre, codepen("css", code));
		break;
	}
    }

    
})(document.getElementsByTagName("pre"),
   // append
   function(pre, node){
       if (pre.nextSibling == null) {
	   pre.parentNode.appendChild(node);
	   return;
       }

       pre.parentNode.insertBefore(node, pre.nextSibling);
   },
   // codepen
   function (key, code, ele, tmp, obj) {
       ele = document.createElement("form");
       ele.target = "codepen";
       ele.method = "POST";
       ele.action = "http://codepen.io/pen/define/";

       tmp = document.createElement("input");
       tmp.name = "data";
       tmp.type = "hidden";
       obj = {};
       obj[key] = code.textContent;
       tmp.value = JSON.stringify(obj);
       ele.appendChild(tmp);

       tmp = document.createElement("button");
       tmp.textContent = "try it!";
       tmp.type = "submit";
       ele.appendChild(tmp);

       return ele
   })
