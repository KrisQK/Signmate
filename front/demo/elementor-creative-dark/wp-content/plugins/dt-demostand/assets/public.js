jQuery(function($) {
    var $demoPanel = $(".demo-panel"),
        $demoSwitch = $(".annoying-button"),
        $contentPanel = $(".content-panel"),
        $closeButton = $(".close-panel"),
        $html = $("html");

    function paintImages() {
        var $thumbnails = $(".load-on-click", $demoPanel)
        $thumbnails.each(function() {
            var $this = $(this);
            $this.attr("src", $this.attr("data-src"));
            $this.attr("srcset", $this.attr("data-srcset"));
        });
        $thumbnails.loaded(function() {
            $(this).parents(".demo-thumb").addClass("thumb-loaded");
        });
    };
    $demoSwitch.one("click", function() {
        paintImages();
    });
    $demoSwitch.on("click", function() {
        if (!$demoPanel.hasClass("act")) {
            $contentPanel.scrollTop(0);
            $demoPanel.addClass("act");
            $html.css({
                "overflow": "hidden",
                "padding-right": window.innerWidth - $html.width()
            });
        } else {
            $demoPanel.removeClass("act");
            $html.css({
                "overflow": "auto",
                "padding-right": 0
            });
        };
        return;
    });
    $closeButton.on("click", function(event) {
        event.preventDefault();
        $demoSwitch.trigger("click");
    });
    var $filterLinks = $(".filter-panel nav a"),
        $demos = $("a", $contentPanel),
        $contentHolder = $(".content-wrap", $contentPanel);
    $filterLinks.on("click", function(event) {
        event.preventDefault();
        var $this = $(this),
            filterClass = $this.attr("data-filter");
        $contentPanel.scrollTop(0);
        $filterLinks.removeClass("act");
        $this.addClass("act");
        $demos.hide();
        $contentHolder.stop().hide();
        $demos.each(function() {
            var $this = $(this);
            if ($this.hasClass(filterClass)) $this.show();
        });
        $contentHolder.fadeIn(300);
    });
});