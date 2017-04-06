(function(cur,x,style){
    Prism.plugins.autoloader.languages_path = '/assets/prism_grammars/';
    document.getElementById("body").addEventListener("keyup", function(e){
	x = e.which || e.keyCode;
	if (x < 37 || x > 40) return;

	if (x == 37 || x == 38) {
	    // left / up
	    cur--;
	    if (cur < 1) cur = 1;
	} else if (x == 39 || x == 40) {
	    // right / down
	    cur++;
	    if (cur > maxPage) cur = maxPage;
	}

	e = document.getElementById("page" + cur + "");
	e.scrollIntoView();
	style = e.currentStyle || window.getComputedStyle(e);
	window.scrollBy(0, -parseInt(style.marginTop)/2);
    });
})(1);
