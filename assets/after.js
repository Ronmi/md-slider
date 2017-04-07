(function(body,cur,x,style){
    Prism.plugins.autoloader.languages_path = '/assets/prism_grammars/';
    body.addEventListener("keyup", function(e){
	x = e.which || e.keyCode;
	if (x < 35 || x > 40) return;
	e.preventDefault();
    })
    body.addEventListener("keydown", function(e){
	x = e.which || e.keyCode;
	if (x < 35 || x > 40) return;
	e.preventDefault();

	switch (x) {
	    case 35:
		// end
		cur = maxPage;
		break;
	    case 36:
		// home
		cur = 1;
		break;
	    case 37:
	    case 38:
		// left / up
		cur--;
		if (cur < 1) cur = 1;
		break;
	    case 39:
	    case 40:
		// right / down
		cur++;
		if (cur > maxPage) cur = maxPage;
		break;
	}

	e = document.getElementById("page" + cur + "");
	e.scrollIntoView();
	style = e.currentStyle || window.getComputedStyle(e);
	window.scrollBy(0, -parseInt(style.marginTop)/2);
    });
})(document.getElementById("body"),1);
