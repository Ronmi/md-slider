/// <reference path="../go.d.ts" />

(function(){
    for (let i = 1; i <= maxPage; i++) {
        const page = document.getElementById("page" + i + "");
        const div = document.createElement("div");
        div.setAttribute("class", "notes");
        div.setAttribute("id", "note" + i + "");
        const pre = document.createElement("pre");
        div.appendChild(pre);

        const next = page.nextSibling;
        if (!next) {
            page.parentNode.appendChild(div);
            continue;
        }
        page.parentNode.insertBefore(div, next);
    }

    if ("undefined" === typeof notes) {
        return;
    }

    for (let i = 1; i <= maxPage; i++) {
        if ("undefined" === typeof notes["p" + i]) {
            continue;
        }

        const div = document.getElementById("note" + i);
        if (!div) {
            continue;
        }

        div.childNodes[0].textContent = decodeURIComponent(notes["p" + i]);
    }

    let slide: Window;
    let cur = 1;

    const btn = document.createElement("button");
    btn.setAttribute("class", "childFrame");
    btn.textContent = "Open slide window";
    btn.addEventListener("click", function (e: MouseEvent) {
        const href = location.href;
        const arr = href.split("?");
        slide = window.open(arr[0], "md-slide");
    });
    document.body.appendChild(btn);

    document.body.addEventListener("keydown", function (e: KeyboardEvent) {
        if (!slide) return;
        const x = e.which || e.keyCode;
        if (x < 35 || x > 40)
            return;
        switch (x) {
        case 35:
            cur = maxPage;
            break;
        case 36:
            cur = 1;
            break;
        case 37:
        case 38:
            cur--;
            if (cur < 1)
                cur = 1;
            break;
        case 39:
        case 40:
            cur++;
            if (cur > maxPage)
                cur = maxPage;
            break;
        }
        slide.scroll(cur);
    });
})();
