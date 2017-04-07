(function(pres, append, i, pre, code, ele, tmp){
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
		ele = document.createElement("form");
		ele.target = "codepen";
		ele.method = "POST";
		ele.action = "http://codepen.io/pen/define/";

		tmp = document.createElement("input");
		tmp.name = "data";
		tmp.type = "hidden";
		tmp.value = JSON.stringify({
		    "html": code.textContent
		});
		ele.appendChild(tmp);

		tmp = document.createElement("button");
		tmp.textContent = "try it!";
		tmp.type = "submit";
		ele.appendChild(tmp);

		append(pre, ele);
		break;
	    case "language-js":
		ele = document.createElement("form");
		ele.target = "codepen";
		ele.method = "POST";
		ele.action = "http://codepen.io/pen/define/";

		tmp = document.createElement("input");
		tmp.name = "data";
		tmp.type = "hidden";
		tmp.value = JSON.stringify({
		    "js": code.textContent
		});
		ele.appendChild(tmp);

		tmp = document.createElement("button");
		tmp.textContent = "try it!";
		tmp.type = "submit";
		ele.appendChild(tmp);

		append(pre, ele);
		break;
	    case "language-css":
		ele = document.createElement("form");
		ele.target = "codepen";
		ele.method = "POST";
		ele.action = "http://codepen.io/pen/define/";

		tmp = document.createElement("input");
		tmp.name = "data";
		tmp.type = "hidden";
		tmp.value = JSON.stringify({
		    "css": code.textContent
		});
		ele.appendChild(tmp);

		tmp = document.createElement("button");
		tmp.textContent = "try it!";
		tmp.type = "submit";
		ele.appendChild(tmp);

		append(pre, ele);
		break;
	}
    }

    
})(document.getElementsByTagName("pre"),
   function(pre, node){
       if (pre.nextSibling == null) {
	   pre.parentNode.appendChild(node);
	   return;
       }

       pre.parentNode.insertBefore(node, pre.nextSibling);
   })
