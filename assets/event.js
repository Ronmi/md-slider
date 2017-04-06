(function(cur,x,style){
    document.getElementById("body").addEventListener("keypress", function(e){
	x = e.which || e.keyCode;
	if (x != 37 && x != 39) return;

	if (x == 37) {
	    // left
	    cur--;
	    if (cur < 1) cur = 1;
	} else if (x == 39) {
	    // right
	    cur++;
	    if (cur > maxPage) cur = maxPage;
	}

	e = document.getElementById("page" + cur + "");
	e.scrollIntoView();
	style = e.currentStyle || window.getComputedStyle(e);
	window.scrollBy(0, -parseInt(style.marginTop)/2);
    });
})(1);
