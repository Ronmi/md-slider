(function() {
    const pres = document.getElementsByTagName("pre");
    const append = function(pre: HTMLPreElement, node: HTMLElement) {
        if (pre.nextSibling == null) {
            pre.parentNode.appendChild(node);
            return;
        }

        pre.parentNode.insertBefore(node, pre.nextSibling);
    };
    const codepen = function(key, code) {
        const ele = document.createElement("form");
        ele.target = "codepen";
        ele.method = "POST";
        ele.action = "http://codepen.io/pen/define/";

        const hid = document.createElement("input");
        hid.name = "data";
        hid.type = "hidden";
        let obj = {};
        obj[key] = code.textContent;
        hid.value = JSON.stringify(obj);
        ele.appendChild(hid);

        const btn = document.createElement("button");
        btn.textContent = "try it!";
        btn.type = "submit";
        ele.appendChild(btn);

        return ele;
    };

    for (let i = 0; i < pres.length; i++) {
        const pre = pres.item(i);

        if (pre.children.length < 1) continue;

        const code = pre.children[0];
        if (code.tagName !== "CODE") continue;

        // 把 codeblock 通通加上 prism plugin 的 class
        if (pre.className !== "") pre.className += " ";
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


})();
