package docs

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(router *gin.Engine) {
	// Serve static api documentation
	router.GET("/", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/html")
		c.String(http.StatusOK, htmlData)
	})
}

// This data was generated from the corresponding README.md section using pandoc:
//     pandoc README.md -o docs.html
const htmlData = `
	<h1 id="api-endpoints-and-examples">API endpoints and examples</h1>
	<h2 id="latest-stats">Latest stats</h2>
	<p>GET https://api.cndy.store/stats/latest?asset_code=CNDY&amp;asset_issuer=GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX</p>
	<div class="sourceCode" id="cb1"><pre class="sourceCode json"><code class="sourceCode json"><a class="sourceLine" id="cb1-1" data-line-number="1"><span class="fu">{</span></a>
	<a class="sourceLine" id="cb1-2" data-line-number="2">  <span class="dt">&quot;status&quot;</span><span class="fu">:</span> <span class="st">&quot;ok&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-3" data-line-number="3">  <span class="dt">&quot;latest&quot;</span><span class="fu">:</span> <span class="fu">{</span></a>
	<a class="sourceLine" id="cb1-4" data-line-number="4">    <span class="dt">&quot;paging_token&quot;</span><span class="fu">:</span> <span class="st">&quot;33825130903777281-1&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-5" data-line-number="5">    <span class="dt">&quot;asset_type&quot;</span><span class="fu">:</span> <span class="st">&quot;credit_alphanum4&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-6" data-line-number="6">    <span class="dt">&quot;asset_code&quot;</span><span class="fu">:</span> <span class="st">&quot;CNDY&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-7" data-line-number="7">    <span class="dt">&quot;asset_issuer&quot;</span><span class="fu">:</span> <span class="st">&quot;GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-8" data-line-number="8">    <span class="dt">&quot;payments&quot;</span><span class="fu">:</span> <span class="dv">4</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-9" data-line-number="9">    <span class="dt">&quot;accounts_with_trustline&quot;</span><span class="fu">:</span> <span class="dv">4</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-10" data-line-number="10">    <span class="dt">&quot;accounts_with_payments&quot;</span><span class="fu">:</span> <span class="dv">2</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-11" data-line-number="11">    <span class="dt">&quot;created_at&quot;</span><span class="fu">:</span> <span class="st">&quot;2018-03-12T18:49:40Z&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-12" data-line-number="12">    <span class="dt">&quot;issued&quot;</span><span class="fu">:</span> <span class="st">&quot;2000.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-13" data-line-number="13">    <span class="dt">&quot;transferred&quot;</span><span class="fu">:</span> <span class="st">&quot;40.0000000&quot;</span></a>
	<a class="sourceLine" id="cb1-14" data-line-number="14">  <span class="fu">}</span></a>
	<a class="sourceLine" id="cb1-15" data-line-number="15"><span class="fu">}</span></a></code></pre></div>
	<h2 id="asset-stats-history">Asset stats history</h2>
	<p>GET https://api.cndy.store/stats?asset_code=CNDY&amp;asset_issuer=GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX[&amp;from=2018-03-03T23:05:40Z&amp;to=2018-03-03T23:05:50Z]</p>
	<p>If not set, <code>from</code> defaults to UNIX timestamp <code>0</code>, <code>to</code> to <code>now</code>.</p>
	<div class="sourceCode" id="cb2"><pre class="sourceCode json"><code class="sourceCode json"><a class="sourceLine" id="cb2-1" data-line-number="1"><span class="fu">{</span></a>
	<a class="sourceLine" id="cb2-2" data-line-number="2">  <span class="dt">&quot;status&quot;</span><span class="fu">:</span> <span class="st">&quot;ok&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-3" data-line-number="3">  <span class="dt">&quot;stats&quot;</span><span class="fu">:</span> <span class="ot">[</span></a>
	<a class="sourceLine" id="cb2-4" data-line-number="4">    <span class="fu">{</span></a>
	<a class="sourceLine" id="cb2-5" data-line-number="5">      <span class="dt">&quot;paging_token&quot;</span><span class="fu">:</span> <span class="st">&quot;33864305300480001-1&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-6" data-line-number="6">      <span class="dt">&quot;asset_type&quot;</span><span class="fu">:</span> <span class="st">&quot;credit_alphanum4&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-7" data-line-number="7">      <span class="dt">&quot;asset_code&quot;</span><span class="fu">:</span> <span class="st">&quot;CNDY&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-8" data-line-number="8">      <span class="dt">&quot;asset_issuer&quot;</span><span class="fu">:</span> <span class="st">&quot;GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-9" data-line-number="9">      <span class="dt">&quot;payments&quot;</span><span class="fu">:</span> <span class="dv">6</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-10" data-line-number="10">      <span class="dt">&quot;accounts_with_trustline&quot;</span><span class="fu">:</span> <span class="dv">5</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-11" data-line-number="11">      <span class="dt">&quot;accounts_with_payments&quot;</span><span class="fu">:</span> <span class="dv">2</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-12" data-line-number="12">      <span class="dt">&quot;created_at&quot;</span><span class="fu">:</span> <span class="st">&quot;2018-03-13T07:29:48Z&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-13" data-line-number="13">      <span class="dt">&quot;issued&quot;</span><span class="fu">:</span> <span class="st">&quot;2000.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-14" data-line-number="14">      <span class="dt">&quot;transferred&quot;</span><span class="fu">:</span> <span class="st">&quot;140.0000000&quot;</span></a>
	<a class="sourceLine" id="cb2-15" data-line-number="15">    <span class="fu">}</span><span class="ot">,</span></a>
	<a class="sourceLine" id="cb2-16" data-line-number="16">    <span class="fu">{</span></a>
	<a class="sourceLine" id="cb2-17" data-line-number="17">      <span class="dt">&quot;paging_token&quot;</span><span class="fu">:</span> <span class="st">&quot;33864305300480001-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-18" data-line-number="18">      <span class="dt">&quot;asset_type&quot;</span><span class="fu">:</span> <span class="st">&quot;credit_alphanum4&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-19" data-line-number="19">      <span class="dt">&quot;asset_code&quot;</span><span class="fu">:</span> <span class="st">&quot;CNDY&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-20" data-line-number="20">      <span class="dt">&quot;asset_issuer&quot;</span><span class="fu">:</span> <span class="st">&quot;GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-21" data-line-number="21">      <span class="dt">&quot;payments&quot;</span><span class="fu">:</span> <span class="dv">6</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-22" data-line-number="22">      <span class="dt">&quot;accounts_with_trustline&quot;</span><span class="fu">:</span> <span class="dv">5</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-23" data-line-number="23">      <span class="dt">&quot;accounts_with_payments&quot;</span><span class="fu">:</span> <span class="dv">2</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-24" data-line-number="24">      <span class="dt">&quot;created_at&quot;</span><span class="fu">:</span> <span class="st">&quot;2018-03-13T07:29:48Z&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-25" data-line-number="25">      <span class="dt">&quot;issued&quot;</span><span class="fu">:</span> <span class="st">&quot;3000.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-26" data-line-number="26">      <span class="dt">&quot;transferred&quot;</span><span class="fu">:</span> <span class="st">&quot;140.0000000&quot;</span></a>
	<a class="sourceLine" id="cb2-27" data-line-number="27">    <span class="fu">}</span></a>
	<a class="sourceLine" id="cb2-28" data-line-number="28">  <span class="ot">]</span></a>
	<a class="sourceLine" id="cb2-29" data-line-number="29"><span class="fu">}</span></a></code></pre></div>
	<h2 id="current-horizon-cursor">Current Horizon cursor</h2>
	<p>GET https://api.cndy.store/stats/cursor</p>
	<div class="sourceCode" id="cb3"><pre class="sourceCode json"><code class="sourceCode json"><a class="sourceLine" id="cb3-1" data-line-number="1"><span class="fu">{</span></a>
	<a class="sourceLine" id="cb3-2" data-line-number="2">  <span class="dt">&quot;status&quot;</span><span class="fu">:</span> <span class="st">&quot;ok&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-3" data-line-number="3">  <span class="dt">&quot;current_cursor&quot;</span><span class="fu">:</span> <span class="st">&quot;33877250331906049-1&quot;</span></a>
	<a class="sourceLine" id="cb3-4" data-line-number="4"><span class="fu">}</span></a></code></pre></div>
	<h2 id="effects">Effects</h2>
	<p>GET https://api.cndy.store/effects?asset_code=CNDY&amp;asset_issuer=GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX[&amp;from=2018-03-03T23:05:40Z&amp;to=2018-03-03T23:05:50Z]</p>
	<p>If not set, <code>from</code> defaults to UNIX timestamp <code>0</code>, <code>to</code> to <code>now</code>.</p>
	<div class="sourceCode" id="cb4"><pre class="sourceCode json"><code class="sourceCode json"><a class="sourceLine" id="cb4-1" data-line-number="1"><span class="fu">{</span></a>
	<a class="sourceLine" id="cb4-2" data-line-number="2">  <span class="dt">&quot;status&quot;</span><span class="fu">:</span> <span class="st">&quot;ok&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-3" data-line-number="3">  <span class="dt">&quot;effects&quot;</span><span class="fu">:</span> <span class="ot">[</span></a>
	<a class="sourceLine" id="cb4-4" data-line-number="4">    <span class="fu">{</span></a>
	<a class="sourceLine" id="cb4-5" data-line-number="5">      <span class="dt">&quot;id&quot;</span><span class="fu">:</span> <span class="st">&quot;0033819672000335873-0000000001&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-6" data-line-number="6">      <span class="dt">&quot;operation&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/operations/33819672000335873&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-7" data-line-number="7">      <span class="dt">&quot;succeeds&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/effects?order=desc&amp;cursor=33819672000335873-1&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-8" data-line-number="8">      <span class="dt">&quot;precedes&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/effects?order=asc&amp;cursor=33819672000335873-1&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-9" data-line-number="9">      <span class="dt">&quot;paging_token&quot;</span><span class="fu">:</span> <span class="st">&quot;33819672000335873-1&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-10" data-line-number="10">      <span class="dt">&quot;account&quot;</span><span class="fu">:</span> <span class="st">&quot;GDNH64DRUT4CY3UJLWQIB655PQ6OG34UGYB4NC5DC4TYWLNJIBCEYTTD&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-11" data-line-number="11">      <span class="dt">&quot;type&quot;</span><span class="fu">:</span> <span class="st">&quot;trustline_created&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-12" data-line-number="12">      <span class="dt">&quot;type_i&quot;</span><span class="fu">:</span> <span class="dv">20</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-13" data-line-number="13">      <span class="dt">&quot;starting_balance&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-14" data-line-number="14">      <span class="dt">&quot;asset_type&quot;</span><span class="fu">:</span> <span class="st">&quot;credit_alphanum4&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-15" data-line-number="15">      <span class="dt">&quot;asset_code&quot;</span><span class="fu">:</span> <span class="st">&quot;CNDY&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-16" data-line-number="16">      <span class="dt">&quot;asset_issuer&quot;</span><span class="fu">:</span> <span class="st">&quot;GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-17" data-line-number="17">      <span class="dt">&quot;signer_public_key&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-18" data-line-number="18">      <span class="dt">&quot;signer_weight&quot;</span><span class="fu">:</span> <span class="dv">0</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-19" data-line-number="19">      <span class="dt">&quot;signer_key&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-20" data-line-number="20">      <span class="dt">&quot;signer_type&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-21" data-line-number="21">      <span class="dt">&quot;created_at&quot;</span><span class="fu">:</span> <span class="st">&quot;2018-03-12T17:03:45Z&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-22" data-line-number="22">      <span class="dt">&quot;amount&quot;</span><span class="fu">:</span> <span class="st">&quot;0.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-23" data-line-number="23">      <span class="dt">&quot;balance&quot;</span><span class="fu">:</span> <span class="st">&quot;0.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-24" data-line-number="24">      <span class="dt">&quot;balance_limit&quot;</span><span class="fu">:</span> <span class="st">&quot;922337203685.4775807&quot;</span></a>
	<a class="sourceLine" id="cb4-25" data-line-number="25">    <span class="fu">}</span><span class="ot">,</span></a>
	<a class="sourceLine" id="cb4-26" data-line-number="26">    <span class="fu">{</span></a>
	<a class="sourceLine" id="cb4-27" data-line-number="27">      <span class="dt">&quot;id&quot;</span><span class="fu">:</span> <span class="st">&quot;0033820110087000065-0000000002&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-28" data-line-number="28">      <span class="dt">&quot;operation&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/operations/33820110087000065&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-29" data-line-number="29">      <span class="dt">&quot;succeeds&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/effects?order=desc&amp;cursor=33820110087000065-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-30" data-line-number="30">      <span class="dt">&quot;precedes&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/effects?order=asc&amp;cursor=33820110087000065-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-31" data-line-number="31">      <span class="dt">&quot;paging_token&quot;</span><span class="fu">:</span> <span class="st">&quot;33820110087000065-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-32" data-line-number="32">      <span class="dt">&quot;account&quot;</span><span class="fu">:</span> <span class="st">&quot;GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-33" data-line-number="33">      <span class="dt">&quot;type&quot;</span><span class="fu">:</span> <span class="st">&quot;account_debited&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-34" data-line-number="34">      <span class="dt">&quot;type_i&quot;</span><span class="fu">:</span> <span class="dv">3</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-35" data-line-number="35">      <span class="dt">&quot;starting_balance&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-36" data-line-number="36">      <span class="dt">&quot;asset_type&quot;</span><span class="fu">:</span> <span class="st">&quot;credit_alphanum4&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-37" data-line-number="37">      <span class="dt">&quot;asset_code&quot;</span><span class="fu">:</span> <span class="st">&quot;CNDY&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-38" data-line-number="38">      <span class="dt">&quot;asset_issuer&quot;</span><span class="fu">:</span> <span class="st">&quot;GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-39" data-line-number="39">      <span class="dt">&quot;signer_public_key&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-40" data-line-number="40">      <span class="dt">&quot;signer_weight&quot;</span><span class="fu">:</span> <span class="dv">0</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-41" data-line-number="41">      <span class="dt">&quot;signer_key&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-42" data-line-number="42">      <span class="dt">&quot;signer_type&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-43" data-line-number="43">      <span class="dt">&quot;created_at&quot;</span><span class="fu">:</span> <span class="st">&quot;2018-03-12T17:12:15Z&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-44" data-line-number="44">      <span class="dt">&quot;amount&quot;</span><span class="fu">:</span> <span class="st">&quot;1000.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-45" data-line-number="45">      <span class="dt">&quot;balance&quot;</span><span class="fu">:</span> <span class="st">&quot;0.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb4-46" data-line-number="46">      <span class="dt">&quot;balance_limit&quot;</span><span class="fu">:</span> <span class="st">&quot;0.0000000&quot;</span></a>
	<a class="sourceLine" id="cb4-47" data-line-number="47">    <span class="fu">}</span></a>
	<a class="sourceLine" id="cb4-48" data-line-number="48"><span class="er">}</span></a></code></pre></div>
	`
