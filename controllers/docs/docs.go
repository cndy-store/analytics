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
	<h1 id="api-documentation">API documentation</h1>
	<h2 id="stats">Stats</h2>
	<p>GET https://api.cndy.store/stats[?from=2018-03-03T23:05:40Z&amp;to=2018-03-03T23:05:50Z]</p>
	<p>If not set, <code>from</code> defaults to UNIX timestamp <code>0</code>, <code>to</code> to <code>now</code>.</p>
	<div class="sourceCode" id="cb1"><pre class="sourceCode json"><code class="sourceCode json"><a class="sourceLine" id="cb1-1" data-line-number="1"><span class="fu">{</span></a>
	<a class="sourceLine" id="cb1-2" data-line-number="2">  <span class="dt">&quot;accounts_involved&quot;</span><span class="fu">:</span> <span class="dv">10</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-3" data-line-number="3">  <span class="dt">&quot;amount_issued&quot;</span><span class="fu">:</span> <span class="st">&quot;4000.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-4" data-line-number="4">  <span class="dt">&quot;amount_transferred&quot;</span><span class="fu">:</span> <span class="st">&quot;6410.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-5" data-line-number="5">  <span class="dt">&quot;asset_code&quot;</span><span class="fu">:</span> <span class="st">&quot;CNDY&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-6" data-line-number="6">  <span class="dt">&quot;current_cursor&quot;</span><span class="fu">:</span> <span class="st">&quot;35496616211263489-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-7" data-line-number="7">  <span class="dt">&quot;effect_count&quot;</span><span class="fu">:</span> <span class="dv">59</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb1-8" data-line-number="8">  <span class="dt">&quot;trustlines_created&quot;</span><span class="fu">:</span> <span class="dv">9</span></a>
	<a class="sourceLine" id="cb1-9" data-line-number="9"><span class="fu">}</span></a></code></pre></div>
	<h2 id="effects">Effects</h2>
	<p>GET https://api.cndy.store/effects[?from=2018-03-03T23:05:40Z&amp;to=2018-03-03T23:05:50Z]</p>
	<p>If not set, <code>from</code> defaults to UNIX timestamp <code>0</code>, <code>to</code> to <code>now</code>.</p>
	<div class="sourceCode" id="cb2"><pre class="sourceCode json"><code class="sourceCode json"><a class="sourceLine" id="cb2-1" data-line-number="1"><span class="fu">{</span></a>
	<a class="sourceLine" id="cb2-2" data-line-number="2">  <span class="dt">&quot;effects&quot;</span><span class="fu">:</span> <span class="ot">[</span></a>
	<a class="sourceLine" id="cb2-3" data-line-number="3">    <span class="fu">{</span></a>
	<a class="sourceLine" id="cb2-4" data-line-number="4">      <span class="dt">&quot;id&quot;</span><span class="fu">:</span> <span class="st">&quot;0034641642841444353-0000000002&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-5" data-line-number="5">      <span class="dt">&quot;operation&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/operations/34641642841444353&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-6" data-line-number="6">      <span class="dt">&quot;succeeds&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/effects?order=desc&amp;cursor=34641642841444353-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-7" data-line-number="7">      <span class="dt">&quot;precedes&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/effects?order=asc&amp;cursor=34641642841444353-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-8" data-line-number="8">      <span class="dt">&quot;paging_token&quot;</span><span class="fu">:</span> <span class="st">&quot;34641642841444353-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-9" data-line-number="9">      <span class="dt">&quot;account&quot;</span><span class="fu">:</span> <span class="st">&quot;GBET4AVL4BYLFJTFKX2UYE3Y3EHNZXOBMBO5FP7O5FFTHSEAPZ5VEHHD&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-10" data-line-number="10">      <span class="dt">&quot;type&quot;</span><span class="fu">:</span> <span class="st">&quot;account_debited&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-11" data-line-number="11">      <span class="dt">&quot;type_i&quot;</span><span class="fu">:</span> <span class="dv">3</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-12" data-line-number="12">      <span class="dt">&quot;starting_balance&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-13" data-line-number="13">      <span class="dt">&quot;asset_type&quot;</span><span class="fu">:</span> <span class="st">&quot;credit_alphanum4&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-14" data-line-number="14">      <span class="dt">&quot;asset_code&quot;</span><span class="fu">:</span> <span class="st">&quot;CNDY&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-15" data-line-number="15">      <span class="dt">&quot;asset_issuer&quot;</span><span class="fu">:</span> <span class="st">&quot;GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-16" data-line-number="16">      <span class="dt">&quot;signer_public_key&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-17" data-line-number="17">      <span class="dt">&quot;signer_weight&quot;</span><span class="fu">:</span> <span class="dv">0</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-18" data-line-number="18">      <span class="dt">&quot;signer_key&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-19" data-line-number="19">      <span class="dt">&quot;signer_type&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-20" data-line-number="20">      <span class="dt">&quot;created_at&quot;</span><span class="fu">:</span> <span class="st">&quot;2018-03-23T18:54:05Z&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-21" data-line-number="21">      <span class="dt">&quot;amount&quot;</span><span class="fu">:</span> <span class="st">&quot;10.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-22" data-line-number="22">      <span class="dt">&quot;balance&quot;</span><span class="fu">:</span> <span class="st">&quot;0.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-23" data-line-number="23">      <span class="dt">&quot;balance_limit&quot;</span><span class="fu">:</span> <span class="st">&quot;0.0000000&quot;</span></a>
	<a class="sourceLine" id="cb2-24" data-line-number="24">    <span class="fu">}</span><span class="ot">,</span></a>
	<a class="sourceLine" id="cb2-25" data-line-number="25">    <span class="fu">{</span></a>
	<a class="sourceLine" id="cb2-26" data-line-number="26">      <span class="dt">&quot;id&quot;</span><span class="fu">:</span> <span class="st">&quot;0034683389923561473-0000000002&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-27" data-line-number="27">      <span class="dt">&quot;operation&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/operations/34683389923561473&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-28" data-line-number="28">      <span class="dt">&quot;succeeds&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/effects?order=desc&amp;cursor=34683389923561473-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-29" data-line-number="29">      <span class="dt">&quot;precedes&quot;</span><span class="fu">:</span> <span class="st">&quot;https://horizon-testnet.stellar.org/effects?order=asc&amp;cursor=34683389923561473-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-30" data-line-number="30">      <span class="dt">&quot;paging_token&quot;</span><span class="fu">:</span> <span class="st">&quot;34683389923561473-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-31" data-line-number="31">      <span class="dt">&quot;account&quot;</span><span class="fu">:</span> <span class="st">&quot;GBET4AVL4BYLFJTFKX2UYE3Y3EHNZXOBMBO5FP7O5FFTHSEAPZ5VEHHD&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-32" data-line-number="32">      <span class="dt">&quot;type&quot;</span><span class="fu">:</span> <span class="st">&quot;account_debited&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-33" data-line-number="33">      <span class="dt">&quot;type_i&quot;</span><span class="fu">:</span> <span class="dv">3</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-34" data-line-number="34">      <span class="dt">&quot;starting_balance&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-35" data-line-number="35">      <span class="dt">&quot;asset_type&quot;</span><span class="fu">:</span> <span class="st">&quot;credit_alphanum4&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-36" data-line-number="36">      <span class="dt">&quot;asset_code&quot;</span><span class="fu">:</span> <span class="st">&quot;CNDY&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-37" data-line-number="37">      <span class="dt">&quot;asset_issuer&quot;</span><span class="fu">:</span> <span class="st">&quot;GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-38" data-line-number="38">      <span class="dt">&quot;signer_public_key&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-39" data-line-number="39">      <span class="dt">&quot;signer_weight&quot;</span><span class="fu">:</span> <span class="dv">0</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-40" data-line-number="40">      <span class="dt">&quot;signer_key&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-41" data-line-number="41">      <span class="dt">&quot;signer_type&quot;</span><span class="fu">:</span> <span class="st">&quot;&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-42" data-line-number="42">      <span class="dt">&quot;created_at&quot;</span><span class="fu">:</span> <span class="st">&quot;2018-03-24T08:24:06Z&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-43" data-line-number="43">      <span class="dt">&quot;amount&quot;</span><span class="fu">:</span> <span class="st">&quot;10.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-44" data-line-number="44">      <span class="dt">&quot;balance&quot;</span><span class="fu">:</span> <span class="st">&quot;0.0000000&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb2-45" data-line-number="45">      <span class="dt">&quot;balance_limit&quot;</span><span class="fu">:</span> <span class="st">&quot;0.0000000&quot;</span></a>
	<a class="sourceLine" id="cb2-46" data-line-number="46">    <span class="fu">}</span></a>
	<a class="sourceLine" id="cb2-47" data-line-number="47">  <span class="ot">]</span></a>
	<a class="sourceLine" id="cb2-48" data-line-number="48"><span class="fu">}</span></a></code></pre></div>
	<h2 id="asset-stats">Asset stats</h2>
	<p>GET https://api.cndy.store/history[?from=2018-03-03T23:05:40Z&amp;to=2018-03-03T23:05:50Z]</p>
	<p>If not set, <code>from</code> defaults to UNIX timestamp <code>0</code>, <code>to</code> to <code>now</code>.</p>
	<div class="sourceCode" id="cb3"><pre class="sourceCode json"><code class="sourceCode json"><a class="sourceLine" id="cb3-1" data-line-number="1"><span class="fu">{</span></a>
	<a class="sourceLine" id="cb3-2" data-line-number="2">  <span class="dt">&quot;history&quot;</span><span class="fu">:</span> <span class="ot">[</span></a>
	<a class="sourceLine" id="cb3-3" data-line-number="3">    <span class="fu">{</span></a>
	<a class="sourceLine" id="cb3-4" data-line-number="4">      <span class="dt">&quot;paging_token&quot;</span><span class="fu">:</span> <span class="st">&quot;34683389923561473-1&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-5" data-line-number="5">      <span class="dt">&quot;asset_type&quot;</span><span class="fu">:</span> <span class="st">&quot;credit_alphanum4&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-6" data-line-number="6">      <span class="dt">&quot;asset_code&quot;</span><span class="fu">:</span> <span class="st">&quot;CNDY&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-7" data-line-number="7">      <span class="dt">&quot;asset_issuer&quot;</span><span class="fu">:</span> <span class="st">&quot;GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-8" data-line-number="8">      <span class="dt">&quot;num_accounts&quot;</span><span class="fu">:</span> <span class="dv">10</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-9" data-line-number="9">      <span class="dt">&quot;payments&quot;</span><span class="fu">:</span> <span class="dv">58</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-10" data-line-number="10">      <span class="dt">&quot;created_at&quot;</span><span class="fu">:</span> <span class="st">&quot;2018-03-24T08:24:06Z&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-11" data-line-number="11">      <span class="dt">&quot;total_amount&quot;</span><span class="fu">:</span> <span class="st">&quot;4000.0000000&quot;</span></a>
	<a class="sourceLine" id="cb3-12" data-line-number="12">    <span class="fu">}</span><span class="ot">,</span></a>
	<a class="sourceLine" id="cb3-13" data-line-number="13">    <span class="fu">{</span></a>
	<a class="sourceLine" id="cb3-14" data-line-number="14">      <span class="dt">&quot;paging_token&quot;</span><span class="fu">:</span> <span class="st">&quot;34683389923561473-2&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-15" data-line-number="15">      <span class="dt">&quot;asset_type&quot;</span><span class="fu">:</span> <span class="st">&quot;credit_alphanum4&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-16" data-line-number="16">      <span class="dt">&quot;asset_code&quot;</span><span class="fu">:</span> <span class="st">&quot;CNDY&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-17" data-line-number="17">      <span class="dt">&quot;asset_issuer&quot;</span><span class="fu">:</span> <span class="st">&quot;GCJKC2MI63KSQ6MLE6GBSXPDKTDAK43WR522ZYR3F34NPM7Z5UEPIZNX&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-18" data-line-number="18">      <span class="dt">&quot;num_accounts&quot;</span><span class="fu">:</span> <span class="dv">10</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-19" data-line-number="19">      <span class="dt">&quot;payments&quot;</span><span class="fu">:</span> <span class="dv">59</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-20" data-line-number="20">      <span class="dt">&quot;created_at&quot;</span><span class="fu">:</span> <span class="st">&quot;2018-03-24T08:24:06Z&quot;</span><span class="fu">,</span></a>
	<a class="sourceLine" id="cb3-21" data-line-number="21">      <span class="dt">&quot;total_amount&quot;</span><span class="fu">:</span> <span class="st">&quot;4000.0000000&quot;</span></a>
	<a class="sourceLine" id="cb3-22" data-line-number="22">    <span class="fu">}</span></a>
	<a class="sourceLine" id="cb3-23" data-line-number="23">  <span class="ot">]</span></a>
	<a class="sourceLine" id="cb3-24" data-line-number="24"><span class="fu">}</span></a></code></pre></div>
	`
