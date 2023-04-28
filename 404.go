package httpreserve

const fourfour = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>httpreserve | 404 | page not found</title>

	<link id="favicon" rel="icon" type="image/x-icon"
	href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAAABm
	JLR0QA/wD/AP+gvaeTAAAACXBIWXMAAAsTAAALEwEAmpwYAAAA
	B3RJTUUH4QMZAiAVgKAUZQAAAGlJREFUWMPt1jkSgDAMQ9GYk/
	vmoqIM4C1qpAPkv5k0NgBYxF2LPAE0rTQADsCZ8WfOjJ9FbOJn
	EB/xWcTP+AwiGO9FJOM9iGK8hmiK5xDN8RhiKL5FWOJrXm9IMw
	u9qZNMAAEEEIAOuAEmTWnhcv1r2AAAAABJRU5ErkJggg=="/>

	<style>
		body { min-width:750px; font-family: arial, verdana; font-size: 10px; margin-bottom: 20px}
		figcaption { font-family: times new roman, arial, verdana; font-size: 20px; font-weight: bold; margin-bottom: 8px; }

		figcaption.loading { font-family: arial, verdana; font-size: 8px; margin-top: -2px; font-weight: normal; }

		div.wrap { margin: 0 auto ; width:715px; }
		div.layout { margin-top: 50px; min-height: 100%; height: 350px; }

		h4 { font-size: 16px; margin-bottom: -2px;}

		input.link { display: block; margin: auto; width: 500px; }
		input.button  { display: block; margin: auto; }

      /*use push to position footer more usefully on screen if necessary*/
      div.push { height: 340px; min-height: 340px; }

      div.footer { height: 50px; margin: 0 auto ; width:200px; text-align: center; }
	</style>

</head>
<body>
<div class="wrap">
	<div class="layout">
	<p>
	<center>
            <figure>
                <figcaption style="font-family: helvetica; arial, verdana; margin-bottom: 8px"><h1>httpreserve</h1></figcaption>
                <img src=" ` + httpreserveImage + `"
                width="80px" height="80px" alt="httpreserve"/>
            </figure>
	</center>
	</p>
	<center>
	<h4>404 | Page Not Found</h4>
	<br/>
	<p>
	I'm sorry, the page you requested cannot be found. It's likely not your fault
	though and I'll do what I can to fix the situation.
	<br/><br/>
	Please visit <a href="https://github.com/httpreserve/httpreserve/issues">GitHub Issues</a>
	to log an issue.
	<br/><br/>
	This project is designed to help us look at web archiving from the inside out
	<br/>
	and to begin fixing the problem of broken links in documentary heritage.
	<br/><br/>
	If you'd like more information on the impact of broken links please watch this
	YouTube video from one of my colleagues in Melbourne.
	<a href="https://www.youtube.com/watch?v=94JyVSFk8-0">Nicola Laurent</a>
	<br>
	at <a href="https://twitter.com/hashtag/asalinks?f=tweets&vertical=default&lang=en">#ASALinks</a>
	<br/><br/>
	Check out the rest of my work on <a href="http://github.com/httpreserve">GitHub.com</a>
	<br/><br/>
	Other background to this work
	<a href="https://www.youtube.com/watch?v=Ked9GRmKlRw">Binary Trees: Automatically Identifying the links between born digital records.</a>
	<br/>
	</p>
	</center>
	<div class="push">&nbsp;</div>
	</div>
   <div class="footer" id="footer">
	A project by <a href="https://twitter.com/beet_keeper" alt="@beet_keeper on Twitter">@beet_keeper</a>
	<br/>
	On GitHub: <a href="https://github.com/exponential-decay/httpreserve" alt="httpreserve on GitHub">httpreserve</a>
	</div>
</div>
</body>
</html>
`
