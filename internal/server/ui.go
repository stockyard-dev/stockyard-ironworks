package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Ironworks</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#5b8dd9;--purple:#9d6bb8;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5;font-size:13px}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.2rem 1.5rem;max-width:1200px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700;color:var(--gold)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.table{background:var(--bg2);border:1px solid var(--bg3);overflow-x:auto}
.table table{width:100%;border-collapse:collapse;font-size:.7rem}
.table th{text-align:left;padding:.6rem .7rem;color:var(--cm);text-transform:uppercase;font-size:.55rem;letter-spacing:1px;border-bottom:1px solid var(--bg3);background:var(--bg)}
.table td{padding:.6rem .7rem;border-bottom:1px solid var(--bg3);color:var(--cream);vertical-align:top}
.table tr:hover td{background:var(--bg3);cursor:pointer}
.table tr.archived td{opacity:.55}
.table tr.broken td{color:var(--red)}
.col-name{font-weight:700}
.col-size{text-align:right;font-family:var(--mono);color:var(--cd)}
.col-version{color:var(--cd)}
.col-platform{color:var(--leather)}
.checksum{font-family:var(--mono);font-size:.55rem;color:var(--cm);overflow:hidden;text-overflow:ellipsis;display:inline-block;max-width:120px;vertical-align:middle}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm);font-weight:700}
.badge.available{border-color:var(--green);color:var(--green)}
.badge.archived{border-color:var(--cm);color:var(--cm)}
.badge.broken{border-color:var(--red);color:var(--red)}
.badge.staging{border-color:var(--orange);color:var(--orange)}

.btn{font-family:var(--mono);font-size:.6rem;padding:.3rem .55rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.btn-del{color:var(--red);border-color:#3a1a1a}
.btn-del:hover{border-color:var(--red);color:var(--red)}

.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:520px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> IRONWORKS</h1>
<button class="btn btn-p" onclick="openNew()">+ Register Artifact</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search name, build_id, version..." oninput="debouncedRender()">
<select class="filter-sel" id="status-filter" onchange="render()">
<option value="">All Statuses</option>
<option value="available">Available</option>
<option value="staging">Staging</option>
<option value="archived">Archived</option>
<option value="broken">Broken</option>
</select>
<select class="filter-sel" id="platform-filter" onchange="render()">
<option value="">All Platforms</option>
</select>
</div>
<div class="table" id="table-wrap"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='artifacts';

var fields=[
{name:'name',label:'Name',type:'text',required:true},
{name:'version',label:'Version',type:'text',placeholder:'e.g. 1.2.3'},
{name:'platform',label:'Platform',type:'select_or_text',options:['linux-amd64','linux-arm64','darwin-amd64','darwin-arm64','windows-amd64']},
{name:'build_id',label:'Build ID',type:'text'},
{name:'size_bytes',label:'Size (bytes)',type:'number'},
{name:'checksum',label:'Checksum',type:'text',placeholder:'sha256:...'},
{name:'status',label:'Status',type:'select',options:['available','staging','archived','broken']},
{name:'download_url',label:'Download URL',type:'url'}
];

var artifacts=[],artExtras={},editId=null,searchTimer=null;

function fmtBytes(n){
if(!n||n<0)return'0 B';
var units=['B','KB','MB','GB','TB'];
var i=0;
while(n>=1024&&i<units.length-1){n/=1024;i++}
return(i===0?n.toFixed(0):n.toFixed(1))+' '+units[i];
}

function fmtDate(s){
if(!s)return'';
try{
var d=new Date(s);
if(isNaN(d.getTime()))return s;
return d.toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'});
}catch(e){return s}
}

function fieldByName(n){
for(var i=0;i<fields.length;i++)if(fields[i].name===n)return fields[i];
return null;
}

function debouncedRender(){
clearTimeout(searchTimer);
searchTimer=setTimeout(render,200);
}

// ─── Loading ──────────────────────────────────────────────────────

async function load(){
try{
var resps=await Promise.all([
fetch(A+'/artifacts').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
artifacts=resps[0].artifacts||[];
renderStats(resps[1]||{});

try{
var ex=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
artExtras=ex||{};
artifacts.forEach(function(a){
var x=artExtras[a.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(a[k]===undefined)a[k]=x[k]});
});
}catch(e){artExtras={}}

populatePlatformFilter();
}catch(e){
console.error('load failed',e);
artifacts=[];
}
render();
}

function populatePlatformFilter(){
var sel=document.getElementById('platform-filter');
if(!sel)return;
var current=sel.value;
var seen={};
var p=[];
artifacts.forEach(function(a){if(a.platform&&!seen[a.platform]){seen[a.platform]=true;p.push(a.platform)}});
p.sort();
sel.innerHTML='<option value="">All Platforms</option>'+p.map(function(x){return'<option value="'+esc(x)+'"'+(x===current?' selected':'')+'>'+esc(x)+'</option>'}).join('');
}

function renderStats(s){
var total=s.total||0;
var totalBytes=s.total_bytes||0;
var distinct=s.distinct_names||0;
var byStatus=s.by_status||{};
var available=byStatus.available||0;
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Artifacts</div></div>'+
'<div class="st"><div class="st-v">'+distinct+'</div><div class="st-l">Distinct Names</div></div>'+
'<div class="st"><div class="st-v">'+available+'</div><div class="st-l">Available</div></div>'+
'<div class="st"><div class="st-v">'+fmtBytes(totalBytes)+'</div><div class="st-l">Total Size</div></div>';
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var sf=document.getElementById('status-filter').value;
var pf=document.getElementById('platform-filter').value;

var f=artifacts.slice();
if(q)f=f.filter(function(a){
return(a.name||'').toLowerCase().includes(q)||
(a.build_id||'').toLowerCase().includes(q)||
(a.version||'').toLowerCase().includes(q);
});
if(sf)f=f.filter(function(a){return a.status===sf});
if(pf)f=f.filter(function(a){return a.platform===pf});

if(!f.length){
var msg=window._emptyMsg||'No artifacts registered yet.';
document.getElementById('table-wrap').innerHTML='<div class="empty">'+esc(msg)+'</div>';
return;
}

var customCols=fields.filter(function(fd){return fd.isCustom});

var h='<table><thead><tr>';
h+='<th>Name</th><th>Version</th><th>Platform</th><th>Status</th><th class="col-size">Size</th><th>Checksum</th>';
customCols.forEach(function(fd){h+='<th>'+esc(fd.label)+'</th>'});
h+='<th>Created</th></tr></thead><tbody>';

f.forEach(function(a){
var cls=(a.status==='archived'?'archived':'')+(a.status==='broken'?'broken':'');
h+='<tr class="'+cls+'" onclick="openEdit(\''+esc(a.id)+'\')">';
h+='<td class="col-name">'+esc(a.name)+'</td>';
h+='<td class="col-version">'+esc(a.version||'-')+'</td>';
h+='<td class="col-platform">'+esc(a.platform||'-')+'</td>';
h+='<td><span class="badge '+esc(a.status||'available')+'">'+esc(a.status||'available')+'</span></td>';
h+='<td class="col-size">'+fmtBytes(a.size_bytes)+'</td>';
h+='<td><span class="checksum" title="'+esc(a.checksum||'')+'">'+esc(a.checksum||'-')+'</span></td>';
customCols.forEach(function(fd){
var v=a[fd.name];
h+='<td>'+(v===undefined||v===null||v===''?'-':esc(String(v)))+'</td>';
});
h+='<td>'+esc(fmtDate(a.created_at))+'</td>';
h+='</tr>';
});
h+='</tbody></table>';

document.getElementById('table-wrap').innerHTML=h;
}

// ─── Modal ────────────────────────────────────────────────────────

function fieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var req=f.required?' *':'';
var ph=f.placeholder?(' placeholder="'+esc(f.placeholder)+'"'):'';
var h='<div class="fr"><label>'+esc(f.label)+req+'</label>';

if(f.type==='select'){
h+='<select id="f-'+f.name+'">';
if(!f.required)h+='<option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=(String(v)===String(o))?' selected':'';
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(String(o))+'</option>';
});
h+='</select>';
}else if(f.type==='select_or_text'){
h+='<input list="dl-'+f.name+'" type="text" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
h+='<datalist id="dl-'+f.name+'">';
(f.options||[]).forEach(function(o){h+='<option value="'+esc(String(o))+'">'});
h+='</datalist>';
}else if(f.type==='textarea'){
h+='<textarea id="f-'+f.name+'" rows="2"'+ph+'>'+esc(String(v))+'</textarea>';
}else if(f.type==='number'){
h+='<input type="number" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}else{
var inputType=f.type||'text';
h+='<input type="'+esc(inputType)+'" id="f-'+f.name+'" value="'+esc(String(v))+'"'+ph+'>';
}
h+='</div>';
return h;
}

function formHTML(art){
var a=art||{};
var isEdit=!!art;
var h='<h2>'+(isEdit?'EDIT ARTIFACT':'NEW ARTIFACT')+'</h2>';

h+=fieldHTML(fieldByName('name'),a.name);
h+='<div class="row2">'+fieldHTML(fieldByName('version'),a.version)+fieldHTML(fieldByName('platform'),a.platform)+'</div>';
h+='<div class="row2">'+fieldHTML(fieldByName('build_id'),a.build_id)+fieldHTML(fieldByName('size_bytes'),a.size_bytes)+'</div>';
h+=fieldHTML(fieldByName('checksum'),a.checksum);
h+='<div class="row2">'+fieldHTML(fieldByName('status'),a.status||'available')+fieldHTML(fieldByName('download_url'),a.download_url)+'</div>';

var customFields=fields.filter(function(f){return f.isCustom});
if(customFields.length){
var label=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
customFields.forEach(function(f){h+=fieldHTML(f,a[f.name])});
h+='</div>';
}

h+='<div class="acts">';
if(isEdit)h+='<button class="btn btn-del" onclick="delItem()">Delete</button>';
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Register')+'</button>';
h+='</div>';
return h;
}

function openNew(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var n=document.getElementById('f-name');
if(n)n.focus();
}

function openEdit(id){
var a=null;
for(var i=0;i<artifacts.length;i++){if(artifacts[i].id===id){a=artifacts[i];break}}
if(!a)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(a);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

async function submit(){
var nameEl=document.getElementById('f-name');
if(!nameEl||!nameEl.value.trim()){alert('Name is required');return}

var body={};
var extras={};
fields.forEach(function(f){
var el=document.getElementById('f-'+f.name);
if(!el)return;
var val;
if(f.type==='number')val=parseFloat(el.value)||0;
else val=el.value.trim();
if(f.isCustom)extras[f.name]=val;
else body[f.name]=val;
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/artifacts/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/artifacts',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Register failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}
closeModal();
load();
}

async function delItem(){
if(!editId)return;
if(!confirm('Delete this artifact?'))return;
await fetch(A+'/artifacts/'+editId,{method:'DELETE'});
closeModal();
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

if(Array.isArray(cfg.platforms)){
var pf=fieldByName('platform');
if(pf)pf.options=cfg.platforms;
}

if(Array.isArray(cfg.custom_fields)){
cfg.custom_fields.forEach(function(cf){
if(!cf||!cf.name||!cf.label)return;
if(fieldByName(cf.name))return;
fields.push({
name:cf.name,
label:cf.label,
type:cf.type||'text',
options:cf.options||[],
isCustom:true
});
});
}
}).catch(function(){
}).finally(function(){
load();
});
})();
</script>
</body>
</html>`
