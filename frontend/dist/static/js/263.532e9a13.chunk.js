"use strict";(self.webpackChunksrc=self.webpackChunksrc||[]).push([[263],{1385:(e,t,s)=>{s.d(t,{A:()=>a});s(9950);const a=s.p+"static/media/xmark-solid.33c889ee6afa3ed57684011608af6238.svg"},25242:(e,t,s)=>{s.d(t,{$:()=>o});var a=s(44414);const o=e=>{const{className:t,label:s,onClick:o}=e;return(0,a.jsx)("button",{className:"w-full text-white font-medium text-sm shadow-lg rounded-lg py-2 ".concat(null!==t&&void 0!==t?t:""),onClick:o,children:s})}},40088:(e,t,s)=>{s.d(t,{t:()=>h});var a=s(9950),o=s(93461),l=s.n(o),r=s(11274),n=s(49027),i=s.n(n),c=s(67818),d=s(44414);const h=(0,a.forwardRef)(((e,t)=>{const s=.65*window.devicePixelRatio,{boost:o,title:n,series:h,height:u,legend:m,tooltip:v,zooming:f,animation:p,lineWidth:b,tickInterval:x,tickPrecision:g,lineColor:y,backgroundColor:w}=e,[_,S]=(0,a.useState)(!1);(0,a.useEffect)((()=>{S((()=>{if(window.WebGLRenderingContext){const e=document.createElement("canvas");return["webgl","experimental-webgl","webgl2","moz-webgl","webkit-3d"].some((t=>{try{return!!e.getContext(t)}catch(s){return!1}}))}return!1})())}),[]);const[C,j]=(0,a.useState)({chart:{zooming:f?{type:"x"}:{},marginTop:20,height:u,animation:p,backgroundColor:w},xAxis:{labels:{style:{color:"#fff"},format:"{value:%H:%M:%S}"},type:"datetime",tickColor:"#fff",lineColor:y},yAxis:{labels:{style:{color:"#fff"},format:g?"{value:".concat(g,"f}"):"{value:0.2f}"},title:{text:""},opposite:!0,lineColor:y,tickInterval:x},tooltip:{enabled:v,followPointer:!0,followTouchMove:!0,xDateFormat:"%Y-%m-%d %H:%M:%S",padding:12},legend:{enabled:m,itemStyle:{color:"#fff"}},plotOptions:{series:{lineWidth:b,turboThreshold:o?10:0,boostThreshold:o?1:0,states:{hover:{enabled:!1}}}},title:{text:n,style:{color:"#fff",fontSize:"10px",fontWeight:"normal"}},boost:{enabled:_,pixelRatio:s},accessibility:{enabled:!1},credits:{enabled:!1},time:{useUTC:!1},series:[h]}),{t:A}=(0,c.Bd)();return(0,a.useEffect)((()=>{i()(l())}),[]),(0,a.useEffect)((()=>{l().setOptions({lang:{resetZoom:A("components.chart.reset_zoom"),resetZoomTitle:A("components.chart.reset_zoom_title")}})}),[A]),(0,a.useEffect)((()=>{j((e=>({...e,chart:{...e.chart,height:u},title:{...e.title,text:n},boost:{...e.boost,enabled:_}})))}),[u,n,_]),(0,d.jsx)(r.HighchartsReact,{ref:t,options:C,highcharts:l()})}))},95660:(e,t,s)=>{s.d(t,{l:()=>p});var a=s(96583),o=s(40033),l=s(28170),r=s(6493),n=s(29925),i=s(25979),c=s(3788),d=s(95383),h=s(44322),u=s(79739),m=s(10226),v=s(9950),f=s(44414);const p=e=>{const{open:t,title:s,content:p,cancelText:b,submitText:x,placeholder:g,defaultValue:y,inputType:w,onSubmit:_,onClose:S,selectOptions:C}=e,j=(0,v.useRef)(null),[A,N]=(0,v.useState)("");return(0,v.useEffect)((()=>{var e;N(null!==(e=null===C||void 0===C?void 0:C[0].value)&&void 0!==e?e:"")}),[C]),(0,f.jsxs)(a.A,{onClose:S,open:t,children:[(0,f.jsx)(o.A,{children:s}),(0,f.jsxs)(l.A,{children:[p&&(0,f.jsx)(r.A,{children:p}),"select"!==w?(0,f.jsx)(n.A,{autoFocus:!0,fullWidth:!0,ref:j,className:"mt-8",type:w,label:g,defaultValue:y}):(0,f.jsxs)(i.A,{sx:{my:2},fullWidth:!0,children:[(0,f.jsx)(c.A,{id:"select",children:g}),(0,f.jsx)(d.A,{labelId:"select",label:g,defaultValue:null===C||void 0===C?void 0:C[0].value,onChange:e=>{let{target:t}=e;return N(null===t||void 0===t?void 0:t.value)},children:null===C||void 0===C?void 0:C.map((e=>{let{value:t,label:s}=e;return(0,f.jsx)(h.A,{value:t,children:s},t)}))})]})]}),(0,f.jsxs)(u.A,{children:[b&&(0,f.jsx)(m.A,{onClick:S,children:b}),(0,f.jsx)(m.A,{onClick:()=>{var e,t;_&&_("select"===w?A:null!==(e=null===(t=j.current)||void 0===t?void 0:t.value)&&void 0!==e?e:"")},children:x})]})]})}},53768:(e,t,s)=>{s.d(t,{h:()=>i,R:()=>c});var a=s(9950);const o=s.p+"static/media/square-caret-up-solid.0573794ec033f5ce25c1076e3ac596e3.svg";const l=s.p+"static/media/ellipsis-solid.e3aced2a80c2b888104322b1ae60b47f.svg";var r=s(1385),n=s(44414);let i=function(e){return e[e.COLLAPSE_DISABLE=0]="COLLAPSE_DISABLE",e[e.COLLAPSE_SHOW=1]="COLLAPSE_SHOW",e[e.COLLAPSE_HIDE=2]="COLLAPSE_HIDE",e}({});const c=e=>{const{label:t,text:s,children:c,collapse:d,advanced:h}=e,[u,m]=(0,a.useState)(!1),[v,f]=(0,a.useState)(!1);(0,a.useEffect)((()=>{const e=d||i.COLLAPSE_DISABLE;m(e===i.COLLAPSE_HIDE)}),[d]);const p=(d||i.COLLAPSE_DISABLE)!==i.COLLAPSE_DISABLE;return(0,n.jsxs)("div",{className:"mb-4 flex flex-col rounded-xl text-gray-700 shadow-lg",children:[(0,n.jsx)("div",{className:"mx-4 rounded-lg overflow-hidden shadow-lg",children:c}),(0,n.jsxs)("div",{className:"p-4",children:[(0,n.jsxs)("h6",{className:"text-md font-bold text-gray-800 flex ".concat(p?"cursor-pointer select-none":""),onClick:()=>p&&m(!u),children:[p&&(0,n.jsx)("img",{className:"mx-1 ".concat(u?"rotate-180":""),src:o,alt:""}),t]}),s&&!u&&(0,n.jsx)("div",{className:"text-md pt-2",children:s.split("\n").map((e=>(0,n.jsx)("div",{children:e},e)))}),h&&(0,n.jsxs)("div",{className:"mt-2 space-y-2 ".concat(u?"hidden":"block"),children:[(0,n.jsxs)("div",{className:"mx-1 cursor-pointer",onClick:()=>f(!v),children:[(0,n.jsx)("img",{className:"size-4 ".concat(v?"hidden":"block"),src:l,alt:""}),(0,n.jsx)("img",{className:"size-4 ".concat(v?"block":"hidden"),src:r.A,alt:""})]}),(0,n.jsx)("div",{className:v?"block":"hidden",children:h})]})]})]})}},67780:(e,t,s)=>{s.d(t,{p:()=>r});var a=s(29925),o=s(94628),l=s(44414);const r=e=>{const{label:t,disabled:s,className:r,defaultValue:n,numberLimit:i,type:c,onValueChange:d}=e,h=(0,o.M)((e=>{let{target:t}=e;if(!d)return;const{value:s}=t;if("number"===c){const e=Number(s);if(isNaN(e))return t.value=n.toString(),void d(n);if(i){const{max:s,min:a}=i;if(e>s||e<a)return t.value=n.toString(),void d(n)}d(e)}else d(s)}),1e3);return(0,l.jsx)(a.A,{size:"small",type:c,label:t,disabled:s,onChange:h,defaultValue:n,className:"w-full ".concat(null!==r&&void 0!==r?r:""),InputLabelProps:{shrink:!0}})}},51682:(e,t,s)=>{s.d(t,{J:()=>o});var a=s(44414);const o=e=>{const{className:t,icon:s,label:o,value:l,unit:r,color:n}=e;return(0,a.jsx)("div",{className:"w-full p-2 ".concat(null!==t&&void 0!==t?t:""),children:(0,a.jsxs)("div",{className:"flex flex-row bg-gradient-to-r rounded-md p-4 shadow-xl ".concat(n?"from-indigo-500 via-purple-500 to-pink-500":"bg-gray-50 hover:bg-gray-100 transition-all"),children:[s&&(0,a.jsx)("img",{className:"bg-white p-2 rounded-md w-8 h-8 md:w-12 md:h-12 self-center",src:s,alt:""}),(0,a.jsxs)("div",{className:"flex flex-col flex-grow ".concat(s?"ml-5":""),children:[(0,a.jsx)("div",{className:"text-sm whitespace-nowrap ".concat(n?"text-gray-50":"text-gray-600"),children:o}),(0,a.jsx)("div",{className:"text-md font-medium flex-nowrap ".concat(n?"text-gray-100":"text-gray-800"),children:"".concat(l," ").concat(null!==r&&void 0!==r?r:"")})]})]})})}},60898:(e,t,s)=>{s.d(t,{Z:()=>o});var a=s(44414);const o=e=>{const{embedded:t,className:s,label:o,sublabel:l,children:r}=e;return(0,a.jsx)("div",{className:"w-full text-gray-800",children:(0,a.jsxs)("div",{className:"flex flex-col shadow-lg rounded-lg",children:[(0,a.jsxs)("div",{className:"px-4 py-3 font-bold",children:[l&&(0,a.jsx)("h6",{className:"text-gray-500 text-xs",children:l}),(0,a.jsx)("h2",{className:t?"text-md":"text-lg",children:o})]}),(0,a.jsx)("div",{className:"p-4 m-2 flex flex-col justify-center gap-4 ".concat(null!==s&&void 0!==s?s:""),children:r})]})})}},2074:(e,t,s)=>{s.d(t,{l:()=>v});var a=s(96583),o=s(88060),l=s(83239),r=s(82053),n=s(10226),i=s(249),c=s(32610),d=s(57357),h=s(83563),u=s(74745),m=s(44414);const v=e=>{const{title:t,open:s,options:v,onClose:f,onSelect:p}=e;return(0,m.jsxs)(a.A,{fullWidth:!0,onClose:f,open:s,children:[(0,m.jsx)(o.A,{className:"bg-violet-500",sx:{position:"relative"},children:(0,m.jsxs)(l.A,{children:[(0,m.jsx)(r.A,{sx:{mt:1,ml:1,flex:1},variant:"h6",children:t}),(0,m.jsx)(n.A,{autoFocus:!0,color:"inherit",onClick:f,children:"X"})]})}),(0,m.jsx)(i.A,{children:!(null===v||void 0===v||!v.length)&&v.map((e=>e.length>1&&(0,m.jsxs)("div",{children:[(0,m.jsx)(c.Ay,{children:(0,m.jsx)(d.A,{onClick:()=>p&&p(e[1]),children:(0,m.jsx)(h.A,{primary:e[0],secondary:e[3===e.length?2:1].split("\n").map((e=>(0,m.jsxs)("span",{children:[e,(0,m.jsx)("br",{})]},e)))})})}),(0,m.jsx)(u.A,{})]},e[1])))})]})}},71359:(e,t,s)=>{s.d(t,{w:()=>o});var a=s(35887);const o=function(e){let t=arguments.length>2&&void 0!==arguments[2]?arguments[2]:2e3;arguments.length>1&&void 0!==arguments[1]&&arguments[1]?a.Ay.error(e,{duration:t}):a.Ay.success(e,{duration:t})}},28418:(e,t,s)=>{s.d(t,{U:()=>a});const a=(e,t)=>{const s=[];for(let a=0;a<e.length;a++)a?s.push((e[a]-e[a-1])/(t/1e3)):s.push(0);return s}},7678:(e,t,s)=>{s.d(t,{W:()=>l,z:()=>o});var a=s(96004);let o=function(e){return e[e.LOW_PASS=0]="LOW_PASS",e[e.HIGH_PASS=1]="HIGH_PASS",e[e.BAND_PASS=2]="BAND_PASS",e}({});const l=(e,t)=>{const{passbandType:s,poles:l,sampleRate:r,lowFreqCorner:n,highFreqCorner:i}=t;let c;c=s===o.LOW_PASS?a.com.oregondsp.signalProcessing.filter.iir.PassbandType.LOWPASS:s===o.HIGH_PASS?a.com.oregondsp.signalProcessing.filter.iir.PassbandType.HIGHPASS:a.com.oregondsp.signalProcessing.filter.iir.PassbandType.BANDPASS;const d=new a.com.oregondsp.signalProcessing.filter.iir.Butterworth(l,c,n,i,1/r),h=new Float32Array(e);return d.filterInPlace(h),Array.from(h)}},76969:(e,t,s)=>{s.d(t,{N:()=>a});const a=(e,t)=>{const s=[];for(let a of e)s.push(a/t);return s}},27156:(e,t,s)=>{s.d(t,{i:()=>a});const a=(e,t,s)=>{const a=2**(t-1),o=[];for(let l of e)o.push(s/a*l);return o}},42105:(e,t,s)=>{s.d(t,{D:()=>o});var a=s(43974);const o=e=>{const t=new Date(e);return(0,a.A)(t,"yyyy-MM-dd HH:mm:ss")}},94628:(e,t,s)=>{s.d(t,{M:()=>a});const a=(e,t)=>{let s=null;return function(){for(var a=arguments.length,o=new Array(a),l=0;l<a;l++)o[l]=arguments[l];null===s&&(s=window.setTimeout((()=>(s=null,e(...o))),t))}}},92263:(e,t,s)=>{s.r(t),s.d(t,{default:()=>G});var a=s(67818),o=s(55297),l=s(60898),r=s(81831),n=s(87403),i=s(43396),c=s(69200),d=s(47775),h=s(65869),u=s(90875),m=s(47315),v=s(64779),f=s(73543),p=s(79299),b=s(40922),x=s(73215),g=s(61721),y=s(69178),w=s(50639),_=s(44414);const S={"zh-CN":{theme:u.g,adapterLocale:b.A},"zh-TW":{theme:m.a,adapterLocale:x.A},"en-US":{theme:v.c,adapterLocale:g.A},"ja-JP":{theme:f.q,adapterLocale:y.A},"ko-KR":{theme:p.U,adapterLocale:w.A}},C=e=>{const{label:t,onChange:s,value:a,defaultValue:o,currentLocale:l}=e,u=Intl.DateTimeFormat().resolvedOptions().timeZone,m=(0,d.A)({},S[l].theme);return(0,_.jsx)(h.A,{theme:m,children:(0,_.jsx)(r.$,{dateAdapter:c.h,adapterLocale:S[l].adapterLocale,children:(0,_.jsx)(n.K,{format:"yyyy-MM-dd HH:mm:ss",className:"w-full",timezone:"system",views:["year","month","day","hours","minutes","seconds"],viewRenderers:{hours:i.Zo,minutes:i.Zo,seconds:i.Zo},onChange:e=>{var t;let a=null!==(t=null===e||void 0===e?void 0:e.valueOf())&&void 0!==t?t:0;isNaN(a)&&(a=0),s(a)},slotProps:{field:{clearable:!0}},label:"".concat(t," - ").concat(u),defaultValue:o,value:null!==a&&void 0!==a?a:0,ampm:!1})})})};var j=s(25242),A=s(53768),N=s(40088),k=s(51682),L=s(9950),E=s(41237),P=s(2074),T=s(82569),H=s(67780),O=s(42074);var D=s(71359),F=s(38526),z=s(19325),R=s(35887);const q=async function(e,t,s,a){if(!(!(arguments.length>4&&void 0!==arguments[4])||arguments[4]))return await R.Ay.promise(e,{loading:t,success:s,error:a});try{return await R.Ay.promise(e,{loading:t,success:s,error:a})}catch{}};var I=s(7678);var W=s(77254),M=s(27156),B=s(76969),Z=s(28418);var U=s(95660);const V=(e,t,s)=>{let{network:a,station:o,location:l}=s;const r=new Date(e);return"".concat(r.getUTCFullYear(),".").concat((e=>{const t=new Date(e.getUTCFullYear(),0,0),s=e.getTime()-t.getTime();return Math.floor(s/864e5)})(r).toString().padStart(3,"0"),".").concat(r.getUTCHours().toString().padStart(2,"0"),".").concat(r.getUTCMinutes().toString().padStart(2,"0"),".").concat(r.getUTCSeconds().toString().padStart(2,"0"),".").concat(r.getUTCMilliseconds().toString().padStart(4,"0"),".").concat(a.slice(0,2),".").concat(o.slice(0,5),".").concat(l.slice(0,2),".").concat(t,".D.sac")};var $=s(42105);const G=e=>{var t,s,r,n,i,c;const{t:d}=(0,a.Bd)(),{station:h}=(0,T.d4)((e=>{let{station:t}=e;return t})),{duration:u}=(0,T.d4)((e=>{let{duration:t}=e;return t})),[m,v]=(0,L.useState)(!h.initialized);(0,L.useEffect)((()=>{v(!h.initialized)}),[h.initialized]);const f=Date.now(),[p,b]=(0,O.ok)(),[x,g]=(0,L.useState)({start:p.has("start")?Number(p.get("start")):f-1e3*u,end:p.has("end")?Number(p.get("end")):f}),y=(e,t)=>g((s=>t?{...s,end:e}:{...s,start:e})),[w,S]=(0,L.useState)({open:!1,inputType:"select"}),[R,G]=(0,L.useState)({open:!1}),J=()=>G({...R,open:!1}),[Y,K]=(0,L.useState)({ehz:{label:"views.history.labels.ehz_detail.label",value:"-"},ehe:{label:"views.history.labels.ehe_detail.label",value:"-"},ehn:{label:"views.history.labels.ehn_detail.label",value:"-"}}),[X,Q]=(0,L.useState)({ehz:{holder:{collapse:A.h.COLLAPSE_HIDE,label:"views.history.charts.ehz.label",text:"views.history.charts.ehz.text"},chart:{buffer:[],backgroundColor:"#d97706",filter:{enabled:!1},ref:(0,L.useRef)(null),series:{name:"EHZ",type:"line",color:"#f1f5f9"}}},ehe:{holder:{collapse:A.h.COLLAPSE_SHOW,label:"views.history.charts.ehe.label",text:"views.history.charts.ehe.text"},chart:{buffer:[],backgroundColor:"#10b981",filter:{enabled:!1},ref:(0,L.useRef)(null),series:{name:"EHE",type:"line",color:"#f1f5f9"}}},ehn:{holder:{collapse:A.h.COLLAPSE_SHOW,label:"views.history.charts.ehn.label",text:"views.history.charts.ehn.text"},chart:{buffer:[],backgroundColor:"#0ea5e9",filter:{enabled:!1},ref:(0,L.useRef)(null),series:{name:"EHN",type:"line",color:"#f1f5f9"}}}}),ee=(e,t,s)=>Q((a=>({...a,[e]:{...a[e],chart:{...a[e].chart,filter:{...a[e].chart.filter,[t?"lowCorner":"highCorner"]:s}}}}))),te=async()=>{const{start:e,end:t}=x;if(!e||!t||e>=t)return void(0,D.w)(d("views.history.toasts.duration_error"),!0);const{backend:s}=z.DR,a={start:e,end:t,channel:"",format:"json"},o=await q((0,F.c)({backend:s,payload:a,timeout:120,throwError:!0,endpoint:z.DR.endpoints.history}),d("views.history.toasts.is_fetching_waveform"),d("views.history.toasts.fetch_waveform_success"),d("views.history.toasts.fetch_waveform_error"));((e,t)=>{if(null===e||void 0===e||!e.data)return;const{adc:s}=W.A.getState().adc,{geophone:a}=W.A.getState().geophone;t((t=>(Object.keys(t).forEach((o=>{if(!e.data.every((e=>o in e)))return;const l=e.data.map((e=>e[o])).map((e=>(0,M.i)(e,s.resolution,s.fullscale))).map((e=>{const t=a.sensitivity/100;return(0,B.N)(e,t)})),r=l.map((e=>{const t=1e3/e.length;return(0,Z.U)(e,t)})),n=l.flat().reduce(((e,t)=>Math.max(Math.abs(e),Math.abs(t))),0),i=r.flat().reduce(((e,t)=>Math.max(Math.abs(e),Math.abs(t))),0);t[o]={...t[o],values:{pgv:n.toFixed(5),pga:i.toFixed(5)},value:"views.history.labels.".concat(o,"_detail.value")}})),t)))})(o,K),((e,t)=>{null!==e&&void 0!==e&&e.data&&t((t=>(Object.keys(t).forEach((s=>{if(!e.data.every((e=>s in e)))return;const a=e.data.map((e=>{let{ts:t,...a}=e;return{data:a[s],ts:t}}));t[s].chart.buffer=a;const{enabled:o,lowCorner:l,highCorner:r}=t[s].chart.filter,{lowFreqCorner:n,highFreqCorner:i}={lowFreqCorner:null!==l&&void 0!==l?l:.1,highFreqCorner:null!==r&&void 0!==r?r:10};t[s].chart={...t[s].chart,title:o?"Band pass [".concat(n,"-").concat(i," Hz]"):""};const c=a.map((e=>{let{ts:t,data:s}=e;const a=o?(0,I.W)(s,{poles:4,lowFreqCorner:n,highFreqCorner:i,sampleRate:s.length,passbandType:I.z.BAND_PASS}):s,l=1e3/a.length;return a.map(((e,s)=>[t-l*(a.length-s),e]))})).reduce(((e,t)=>e.concat(t)),[]),{current:d}=t[s].chart.ref;if(d){const{series:e}=d.chart;e[0].setData(c,!0,!1,!1)}})),t)))})(o,Q)},{locale:se}=e,{fallback:ae}=E.X;return(0,_.jsxs)(_.Fragment,{children:[(0,_.jsxs)(o.m,{className:"my-6 gap-4 grid md:grid-cols-2 ".concat(m?"cursor-progress":""),children:[(0,_.jsxs)(l.Z,{label:d("views.history.panels.query_history"),children:[(0,_.jsx)(C,{value:x.start,currentLocale:null!==se&&void 0!==se?se:ae,label:d("views.history.time_pickers.start_time"),onChange:e=>y(e,!1)}),(0,_.jsx)(C,{value:x.end,currentLocale:null!==se&&void 0!==se?se:ae,label:d("views.history.time_pickers.end_time"),onChange:e=>y(e,!0)}),(0,_.jsx)(j.$,{className:"bg-indigo-700 hover:bg-indigo-800 ".concat(m?"cursor-wait":""),onClick:async()=>{m||(v(!0),await te(),v(!1))},label:d("views.history.buttons.query_waveform")}),(0,_.jsx)(j.$,{className:"bg-green-700 hover:bg-green-800",onClick:()=>{const{start:e,end:t}=x;if(!e||!t||e>=t)return void(0,D.w)(d("views.history.toasts.duration_error"),!0);const s=async s=>{S((e=>({...e,open:!1})));const{backend:a}=z.DR,o={start:e,end:t,channel:s,format:"sac"},l=V(e,s,h);await q((0,F.c)({backend:a,payload:o,timeout:120,throwError:!0,endpoint:z.DR.endpoints.history,blobOptions:{filename:l}}),d("views.history.toasts.is_exporting_sac"),d("views.history.toasts.export_sac_success"),d("views.history.toasts.export_sac_error"))};S((e=>({...e,open:!0,selectOptions:[{label:"EHZ",value:"EHZ"},{label:"EHE",value:"EHE"},{label:"EHN",value:"EHN"}],onSubmit:s,title:"views.history.forms.choose_channel.title",cancelText:"views.history.forms.choose_channel.cancel",submitText:"views.history.forms.choose_channel.submit",placeholder:"views.history.forms.choose_channel.placeholder"})))},label:d("views.history.buttons.query_sac_file")}),(0,_.jsx)(j.$,{className:"bg-yellow-700 hover:bg-yellow-800 ".concat(m?"cursor-wait":""),onClick:async()=>{m||(v(!0),await(async()=>{const{backend:e}=z.DR,t=await q((0,F.c)({backend:e,payload:{source:"show"},timeout:30,throwError:!0,endpoint:z.DR.endpoints.trace}),d("views.history.toasts.is_fetching_source"),d("views.history.toasts.fetch_source_success"),d("views.history.toasts.fetch_source_error"));if(null===t||void 0===t||!t.data)return;const s=async t=>{S((e=>({...e,open:!1})));const s=await q((0,F.c)({backend:e,timeout:60,throwError:!0,payload:{source:t},endpoint:z.DR.endpoints.trace}),d("views.history.toasts.is_fetching_events"),d("views.history.toasts.fetch_events_success"),d("views.history.toasts.fetch_events_error"));if(null===s||void 0===s||!s.data)return;const a=e=>{G((e=>({...e,open:!1})));const[t,s]=e.split("|").map(Number);g({start:t,end:s}),(0,D.w)(d("views.history.toasts.event_select_success"))},o=s.data.map((e=>{let{distance:t,magnitude:s,region:a,event:o,timestamp:l,depth:r,estimation:n}=e;return[a,"".concat(l+1e3*n.p,"|").concat(l+1e3*n.s),d("views.history.selects.choose_event.template",{event:o,time:(0,$.D)(l),magnitude:s.toFixed(1),distance:t.toFixed(1),p_wave:n.p.toFixed(1),s_wave:n.s.toFixed(1),depth:-1!==r?r.toFixed(1):"-"})]}));G((e=>({...e,open:!0,options:o,onClose:J,onSelect:a,title:"views.history.selects.choose_event.title"})))};S((e=>({...e,open:!0,selectOptions:t.data.map((e=>"name"in e&&"value"in e?{label:e.name,value:e.value}:{label:"",value:""})),onSubmit:s,title:"views.history.forms.choose_source.title",cancelText:"views.history.forms.choose_source.cancel",submitText:"views.history.forms.choose_source.submit",placeholder:"views.history.forms.choose_source.placeholder"})))})(),v(!1))},label:d("views.history.buttons.query_source")}),(0,_.jsx)(j.$,{className:"bg-cyan-700 hover:bg-cyan-800",onClick:async()=>{const{start:e,end:t}=x;if(!e||!t||e>=t)return void(0,D.w)(d("views.history.toasts.duration_error"),!0);const s=new URLSearchParams;s.set("start",String(e)),s.set("end",String(t)),b(s);const a=window.location.href,o=await(async e=>{var t;const s=null!==(t=navigator.clipboard)&&void 0!==t?t:{writeText:e=>{const t=document.createElement("input");t.value=e,document.body.appendChild(t),t.select(),document.execCommand("copy"),document.body.removeChild(t)}};return!!s&&(s.writeText(e),!0)})(a);(0,D.w)(d(o?"views.history.toasts.copy_link_success":"views.history.toasts.copy_link_error"),!o)},label:d("views.history.buttons.get_share_link")})]}),(0,_.jsx)(l.Z,{className:"",label:d("views.history.panels.analyze_history"),children:Object.values(Y).map((e=>{let{label:t,value:s,values:a,...o}=e;return(0,L.createElement)(k.J,{...o,key:t,label:d(t),value:d(s,a)})}))}),(0,_.jsx)(U.l,{...w,onClose:()=>S({...w,open:!1}),title:d(null!==(t=w.title)&&void 0!==t?t:""),cancelText:d(null!==(s=w.cancelText)&&void 0!==s?s:""),submitText:d(null!==(r=w.submitText)&&void 0!==r?r:""),placeholder:d(null!==(n=w.placeholder)&&void 0!==n?n:""),content:d(null!==(i=w.content)&&void 0!==i?i:"",{...w.values})}),(0,_.jsx)(P.l,{...R,onClose:J,title:d(null!==(c=R.title)&&void 0!==c?c:"")})]}),Object.keys(X).map((e=>{var t,s;return(0,_.jsx)(A.R,{text:d(null!==(t=X[e].holder.text)&&void 0!==t?t:""),label:d(null!==(s=X[e].holder.label)&&void 0!==s?s:""),advanced:(0,_.jsx)(o.m,{className:"max-w-96",children:(0,_.jsxs)(l.Z,{label:d("views.history.charts.".concat(e,".advanced.panels.butterworth_filter.title")),embedded:!0,children:[(0,_.jsxs)(o.m,{className:"flex flex-col md:flex-row gap-4",children:[(0,_.jsx)(H.p,{onValueChange:t=>ee(e,!0,Number(t)),defaultValue:.1,type:"number",disabled:X[e].chart.filter.enabled,numberLimit:{max:100,min:.1},label:d("views.history.charts.".concat(e,".advanced.panels.butterworth_filter.low_corner_freq"))}),(0,_.jsx)(H.p,{onValueChange:t=>ee(e,!1,Number(t)),defaultValue:10,type:"number",disabled:X[e].chart.filter.enabled,numberLimit:{max:100,min:.1},label:d("views.history.charts.".concat(e,".advanced.panels.butterworth_filter.high_corner_freq"))})]}),(0,_.jsx)(j.$,{label:d("views.history.charts.".concat(e,".advanced.panels.butterworth_filter.").concat(X[e].chart.filter.enabled?"disable_filter":"enable_filter")),className:"bg-indigo-600 hover:bg-indigo-700",onClick:()=>{return t=e,void Q((e=>{const s=!e[t].chart.filter.enabled,{lowCorner:a,highCorner:o}=e[t].chart.filter,{lowFreqCorner:l,highFreqCorner:r}={lowFreqCorner:null!==a&&void 0!==a?a:.1,highFreqCorner:null!==o&&void 0!==o?o:10},n=e[t].chart.buffer.map((e=>{let{ts:t,data:a}=e;const o=s?(0,I.W)(a,{poles:4,lowFreqCorner:l,highFreqCorner:r,sampleRate:a.length,passbandType:I.z.BAND_PASS}):a,n=1e3/o.length;return o.map(((e,s)=>[t-n*(o.length-s),e]))})).reduce(((e,t)=>e.concat(t)),[]),{current:i}=e[t].chart.ref;if(i){const{series:e}=i.chart;e[0].setData(n,!0,!1,!1)}const c={...e[t],chart:{...e[t].chart,filter:{...e[t].chart.filter,enabled:s},title:s?"Band pass [".concat(l,"-").concat(r," Hz]"):""}};return{...e,[t]:c}}));var t}})]})}),children:(0,_.jsx)(N.t,{...X[e].chart,height:300,boost:!0,lineWidth:1,tooltip:!0,zooming:!0,animation:!0,tickPrecision:1,tickInterval:100})},X[e].holder.label)}))]})}}}]);