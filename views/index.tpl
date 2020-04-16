<!DOCTYPE HTML>
<html lang="en">
<head>
<!--[if IE]>
<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
<![endif]-->
<meta charset="utf-8">
<title>我们的相册 - {{.Title}}</title>
<meta name="description" content="blueimp Gallery is a touch-enabled, responsive and customizable image and video gallery, carousel and lightbox, optimized for both mobile and desktop web browsers. It features swipe, mouse and keyboard navigation, transition effects, slideshow functionality, fullscreen support and on-demand content loading and can be extended to display additional content types.">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="/css/blueimp-gallery.css">
<link rel="stylesheet" href="/css/blueimp-gallery-indicator.css">
<link rel="stylesheet" href="/css/blueimp-gallery-video.css">
<link rel="stylesheet" href="/css/demo/demo.css">
<link href="http://vjs.zencdn.net/4.12/video-js.css" rel="stylesheet">
<script src="http://vjs.zencdn.net/4.12/video.js"></script>
</head>
<body>
<h1>{{.Title}}</h1>
<p>这里展示了一些我们的私人照片，自己看看就好，切勿随意传播。</p>
<p>点击小图片放大显示。</p>
<!-- The container for the list of example images -->
<div id="links" class="links">
    {{range .Photos}}
    <a href="{{.Big}}" title="{{.Title}}" data-gallery=""><img src="{{.Small}}"/></a>
    {{end}}
</div>
<!-- The Gallery as lightbox dialog, should be a child element of the document body -->
<div id="blueimp-gallery" class="blueimp-gallery">
    <div class="slides"></div>
    <h3 class="title"></h3>
    <a class="prev">‹</a>
    <a class="next">›</a>
    <a class="close">×</a>
    <a class="play-pause"></a>
    <ol class="indicator"></ol>
</div>
<div id="videos" class="videos">
    {{if not (or .IsMobile .IsTablet)}}
    {{range .Videos}}
    <p>{{.Title}}</p>    
    <video width="800" height="600" class="video-js vjs-default-skin" controls preload="auto" data-setup="{}">
        <source src="{{.Url}}" type="video/mp4">
        <p class="vjs-no-js">To view this video please enable JavaScript, and consider upgrading to a web browser that <a href="http://videojs.com/html5-video-support/" target="_blank">supports HTML5 video</a></p>
    </video>
    {{end}}
    {{else}}
    {{range .Videos}}
    <p>{{.Title}}</p>    
    <video width="480" height="320" class="video-js vjs-default-skin" controls preload="auto" data-setup="{}">
        <source src="{{.Url}}" type="video/mp4">
        <p class="vjs-no-js">To view this video please enable JavaScript, and consider upgrading to a web browser that <a href="http://videojs.com/html5-video-support/" target="_blank">supports HTML5 video</a></p>
    </video>
    {{end}}
    {{end}}
</div>
<ul class="navigation">
    {{range .Links}}
	<li><a href="{{.Url}}">{{.Title}}</a></li>
    {{end}}
</ul>
<script src="/js/blueimp-helper.js"></script>
<script src="/js/blueimp-gallery.js"></script>
<script src="/js/blueimp-gallery-fullscreen.js"></script>
<script src="/js/blueimp-gallery-indicator.js"></script>
<script src="/js/blueimp-gallery-video.js"></script>
<script src="/js/blueimp-gallery-vimeo.js"></script>
<script src="/js/blueimp-gallery-youtube.js"></script>
<script src="/js/vendor/jquery.js"></script>
<script src="/js/jquery.blueimp-gallery.js"></script>
</body>
</html>
