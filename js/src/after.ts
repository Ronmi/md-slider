/// <reference path="../go.d.ts" />



(function () {
    const scroll = function (cur: number) {
        const e = document.getElementById("page" + cur + "");
        e.scrollIntoView();
        const style = window.getComputedStyle(e);
        window.scrollBy(0, -parseInt(style.marginTop) / 2);
    };
    const body = document.getElementsByTagName("body")[0];
    let cur = 1;

    // enable prismjs
    if (typeof Prism !== 'undefined') {
        // @ts-ignore
        document.querySelectorAll('code').forEach((el) => {
            Prism.highlightElement(el);
        });
    }

  
    body.addEventListener("keyup", function(e: KeyboardEvent) {
        const x = e.which || e.keyCode;
        if (x < 35 || x > 40) return;
        e.preventDefault();
    });
    body.addEventListener("keydown", function(e: KeyboardEvent) {
        if (e.getModifierState("Fn") || e.getModifierState("Hyper") || e.getModifierState("OS") || e.getModifierState("Super") || e.getModifierState("Win")) return;
        if (e.getModifierState("Control") || e.getModifierState("Alt") || e.getModifierState("Meta") || e.getModifierState("Shift")) return;

        const x = e.which || e.keyCode;
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

        scroll(cur);
    });

    body.onresize = function() {
        scroll(cur);
    };

    scroll(cur);

})();
