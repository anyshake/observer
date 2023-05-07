"use strict";(self.webpackChunksrc=self.webpackChunksrc||[]).push([[460,526,887,874,13,395,481,635,529],{2460:function(e,t,s){s.r(t),s.d(t,{default:function(){return j}});var a=s(3433),l=s(1413),r=s(5671),n=s(3144),i=s(136),o=s(516),c=s(7313),d=s(8395),h=s(9526),x=s(6013),f=s(3481),m=s(9635),u=s(6287),p=s(8446),v=s(2529),w=s(3874),g=s(6887),b=s(6417),j=function(e){(0,i.Z)(s,e);var t=(0,o.Z)(s);function s(e){var a;return(0,r.Z)(this,s),(a=t.call(this,e)).analyseData=function(e){var t=e.acceleration;a.setState({analysis:{vertical:t.vertical[t.vertical.length-1],east_west:t.east_west[t.east_west.length-1],north_south:t.north_south[t.north_south.length-1],synthesis:t.synthesis[t.synthesis.length-1]},lastDataArr:t})},a.state={sidebarMark:"waveform",webSocket:null,lastDataArr:{timestamp:-1,vertical:[],east_west:[],north_south:[]},waveform:{factors:[{name:"\u5782\u76f4\u5206\u91cf",color:"#d97706",data:[]},{name:"\u6c34\u5e73 EW",color:"#0d9488",data:[]},{name:"\u6c34\u5e73 NS",color:"#4f46e5",data:[]}],synthesis:[{name:"\u5408\u6210\u5206\u91cf",color:"#be185d",data:[]}],options:{stroke:{width:2,curve:"smooth"},hollow:{margin:15,size:"40%"},chart:{height:350,toolbar:{show:!1},zoom:{enabled:!1},animations:{enabled:!1}},dataLabels:{enabled:!1},legend:{show:!1,labels:{useSeriesColors:!0}},tooltip:{enabled:!0,theme:"dark",fillSeriesColor:!1,x:{format:"20yy/MM/dd HH:mm:ss"}},xaxis:{type:"datetime",labels:{datetimeUTC:!1,datetimeFormatter:{hour:"HH:mm:ss"},style:{colors:"#fff"}}},yaxis:{opposite:!0,labels:{style:{colors:"#fff"}}}}},response:{uuid:"",station:"",acceleration:[{timestamp:-1,altitude:-1,latitude:-1,longitude:-1,vertical:[],east_west:[],north_south:[]}]},analysis:{vertical:0,east_west:0,north_south:0,synthesis:0}},a}return(0,n.Z)(s,[{key:"componentDidMount",value:function(){var e=this,t=(0,m.default)({tls:u.default.backend.tls,host:u.default.backend.host,port:u.default.backend.port,version:u.default.backend.version,api:u.default.backend.api.socket.uri,type:u.default.backend.api.socket.type});this.setState({webSocket:(0,f.default)({url:t,type:u.default.backend.api.socket.method,onMessageCallback:function(t){var s=t.data,a=JSON.parse(s);e.setState({response:a}),e.drawWaveform(a),e.analyseData(a)}})})}},{key:"componentWillUnmount",value:function(){this.state.webSocket&&(this.state.webSocket.close(),this.setState({webSocket:null}))}},{key:"fillData",value:function(e,t){if(this.state.lastDataArr.timestamp<1||this.state.lastDataArr.timestamp===e.timestamp)return[];var s=(e.timestamp-this.state.lastDataArr.timestamp)/e[t].length;return e[t].map((function(a,l){return[new Date(e.timestamp-(e[t].length-l)*s),a]}))}},{key:"drawWaveform",value:function(e){var t=this,s=e.acceleration;this.state.waveform.synthesis[0].data.length>10500&&this.state.waveform.synthesis[0].data.splice(0,35),this.state.waveform.factors.forEach((function(e,s){t.state.waveform.factors[s].data.length>10500&&t.state.waveform.factors[s].data.splice(0,35)})),this.setState({waveform:(0,l.Z)((0,l.Z)({},this.state.waveform),{},{factors:[(0,l.Z)((0,l.Z)({},this.state.waveform.factors[0]),{},{data:[].concat((0,a.Z)(this.state.waveform.factors[0].data),(0,a.Z)(this.fillData(s,"vertical")))}),(0,l.Z)((0,l.Z)({},this.state.waveform.factors[1]),{},{data:[].concat((0,a.Z)(this.state.waveform.factors[1].data),(0,a.Z)(this.fillData(s,"east_west")))}),(0,l.Z)((0,l.Z)({},this.state.waveform.factors[2]),{},{data:[].concat((0,a.Z)(this.state.waveform.factors[2].data),(0,a.Z)(this.fillData(s,"north_south")))})],synthesis:[(0,l.Z)((0,l.Z)({},this.state.waveform.synthesis[0]),{},{data:[].concat((0,a.Z)(this.state.waveform.synthesis[0].data),(0,a.Z)(this.fillData(s,"synthesis")))})]})})}},{key:"render",value:function(){return(0,b.jsxs)(b.Fragment,{children:[(0,b.jsx)(d.default,{sidebarMark:this.state.sidebarMark}),(0,b.jsxs)("div",{className:"content ml-12 transform ease-in-out duration-500 pt-20 px-2 md:px-5 pb-4",children:[(0,b.jsx)(g.default,{navigation:"\u5b9e\u65f6\u6ce2\u5f62"}),(0,b.jsx)(w.default,{className:0!==this.state.response.uuid.length?"shadow-xl p-4 mb-4 text-sm text-white rounded-lg bg-gradient-to-r from-cyan-500 to-yellow-500":"shadow-xl p-4 mb-4 text-sm text-white rounded-lg bg-gradient-to-r from-blue-500 to-orange-500",icon:0!==this.state.response.uuid.length?(0,b.jsx)("svg",{xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 448 512",className:"w-6 h-6 ml-3",fill:"currentColor",children:(0,b.jsx)("path",{d:"M0 64C0 46.3 14.3 32 32 32c229.8 0 416 186.2 416 416c0 17.7-14.3 32-32 32s-32-14.3-32-32C384 253.6 226.4 96 32 96C14.3 96 0 81.7 0 64zM0 416a64 64 0 1 1 128 0A64 64 0 1 1 0 416zM32 160c159.1 0 288 128.9 288 288c0 17.7-14.3 32-32 32s-32-14.3-32-32c0-123.7-100.3-224-224-224c-17.7 0-32-14.3-32-32s14.3-32 32-32z"})}):(0,b.jsx)("svg",{xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 640 512",className:"w-6 h-6 ml-3",fill:"currentColor",children:(0,b.jsx)("path",{d:"M579.8 267.7c56.5-56.5 56.5-148 0-204.5c-50-50-128.8-56.5-186.3-15.4l-1.6 1.1c-14.4 10.3-17.7 30.3-7.4 44.6s30.3 17.7 44.6 7.4l1.6-1.1c32.1-22.9 76-19.3 103.8 8.6c31.5 31.5 31.5 82.5 0 114L422.3 334.8c-31.5 31.5-82.5 31.5-114 0c-27.9-27.9-31.5-71.8-8.6-103.8l1.1-1.6c10.3-14.4 6.9-34.4-7.4-44.6s-34.4-6.9-44.6 7.4l-1.1 1.6C206.5 251.2 213 330 263 380c56.5 56.5 148 56.5 204.5 0L579.8 267.7zM60.2 244.3c-56.5 56.5-56.5 148 0 204.5c50 50 128.8 56.5 186.3 15.4l1.6-1.1c14.4-10.3 17.7-30.3 7.4-44.6s-30.3-17.7-44.6-7.4l-1.6 1.1c-32.1 22.9-76 19.3-103.8-8.6C74 372 74 321 105.5 289.5L217.7 177.2c31.5-31.5 82.5-31.5 114 0c27.9 27.9 31.5 71.8 8.6 103.9l-1.1 1.6c-10.3 14.4-6.9 34.4 7.4 44.6s34.4 6.9 44.6-7.4l1.1-1.6C433.5 260.8 427 182 377 132c-56.5-56.5-148-56.5-204.5 0L60.2 244.3z"})}),title:0!==this.state.response.uuid.length?"\u6700\u540e\u66f4\u65b0\u4e8e ".concat((0,v.default)(new Date(this.state.response.acceleration.timestamp))):"\u6682\u672a\u6536\u5230\u6570\u636e",text:0!==this.state.response.uuid.length?"".concat(this.state.response.station," - ").concat(this.state.response.uuid):"\u6b63\u5728\u7b49\u5f85\u670d\u52a1\u5668\u6570\u636e..."}),(0,b.jsxs)("div",{className:"flex flex-wrap mt-6",children:[(0,b.jsx)("div",{className:"w-full mb-12 xl:mb-0 px-4",children:(0,b.jsxs)("div",{className:"relative flex flex-col w-full mb-6 shadow-lg rounded-lg",children:[(0,b.jsx)("div",{className:"px-4 py-3  bg-transparent",children:(0,b.jsx)("div",{className:"flex flex-wrap items-center",children:(0,b.jsxs)("div",{className:"relative w-full max-w-full flex-grow flex-1",children:[(0,b.jsx)("h6",{className:"text-gray-500 mb-1 text-xs font-semibold",children:"\u5373\u65f6"}),(0,b.jsx)("h2",{className:"text-gray-700 text-xl font-semibold",children:"\u5b9e\u65f6\u5206\u91cf\u52a0\u901f\u5ea6"})]})})}),(0,b.jsx)("div",{className:"p-4 flex-auto shadow-lg bg-gradient-to-tr from-purple-300 to-purple-400 shadow-purple-500/40 rounded-lg",children:(0,b.jsx)("div",{className:"relative h-[350px]",children:(0,b.jsx)(p.Z,{height:"350px",series:this.state.waveform.factors,options:this.state.waveform.options})})})]})}),(0,b.jsx)("div",{className:"w-full mb-12 xl:mb-0 px-4",children:(0,b.jsxs)("div",{className:"relative flex flex-col w-full mb-6 shadow-lg rounded-lg",children:[(0,b.jsx)("div",{className:"px-4 py-3 bg-transparent",children:(0,b.jsx)("div",{className:"flex flex-wrap items-center",children:(0,b.jsxs)("div",{className:"relative w-full max-w-full flex-grow flex-1",children:[(0,b.jsx)("h6",{className:"text-gray-500 mb-1 text-xs font-semibold",children:"\u5373\u65f6"}),(0,b.jsx)("h2",{className:"text-gray-700 text-xl font-semibold",children:"\u5b9e\u65f6\u5408\u6210\u52a0\u901f\u5ea6"})]})})}),(0,b.jsx)("div",{className:"p-4 flex-auto shadow-lg bg-gradient-to-tr from-indigo-300 to-indigo-400 shadow-indigo-500/40 rounded-lg",children:(0,b.jsx)("div",{className:"relative h-[350px]",children:(0,b.jsx)(p.Z,{type:"area",height:"350px",series:this.state.waveform.synthesis,options:this.state.waveform.options})})})]})}),(0,b.jsx)("div",{className:"w-full px-4",children:(0,b.jsxs)("div",{className:"relative flex flex-col bg-white w-full mb-6 shadow-lg rounded-lg",children:[(0,b.jsx)("div",{className:"px-4 py-3 bg-transparent",children:(0,b.jsx)("div",{className:"flex flex-wrap items-center",children:(0,b.jsxs)("div",{className:"relative w-full max-w-full flex-grow flex-1",children:[(0,b.jsx)("h6",{className:"text-gray-500 mb-1 text-xs font-semibold",children:"\u6570\u636e"}),(0,b.jsx)("h2",{className:"text-gray-700 text-xl font-semibold",children:"\u6570\u636e\u5206\u6790"})]})})}),(0,b.jsx)("div",{className:"p-4 shadow-lg flex-auto",children:(0,b.jsx)("div",{className:"relative h-[350px]",children:(0,b.jsx)("div",{className:"flex flex-wrap -mx-2",children:(0,b.jsxs)("div",{className:"w-full px-2",children:[(0,b.jsx)("div",{className:"relative flex flex-col min-w-0 break-words bg-sky-100 w-full mb-4 shadow-lg rounded-lg",children:(0,b.jsx)("div",{className:"px-4 py-3 bg-transparent",children:(0,b.jsx)("div",{className:"flex flex-wrap items-center",children:(0,b.jsxs)("div",{className:"relative w-full max-w-full flex-grow flex-1",children:[(0,b.jsx)("h6",{className:"text-gray-500 mb-1 text-xs font-semibold",children:"\u5782\u76f4\u5206\u91cf\u5f53\u524d\u503c"}),(0,b.jsx)("h2",{className:"text-gray-700 text-xl font-semibold",children:this.state.analysis.vertical})]})})})}),(0,b.jsx)("div",{className:"relative flex flex-col min-w-0 break-words bg-sky-100 w-full mb-4 shadow-lg rounded-lg",children:(0,b.jsx)("div",{className:"px-4 py-3 bg-transparent",children:(0,b.jsx)("div",{className:"flex flex-wrap items-center",children:(0,b.jsxs)("div",{className:"relative w-full max-w-full flex-grow flex-1",children:[(0,b.jsx)("h6",{className:"text-gray-500 mb-1 text-xs font-semibold",children:"EW \u5206\u91cf\u5f53\u524d\u503c"}),(0,b.jsx)("h2",{className:"text-gray-700 text-xl font-semibold",children:this.state.analysis.east_west})]})})})}),(0,b.jsx)("div",{className:"relative flex flex-col min-w-0 break-words bg-sky-100 w-full mb-4 shadow-lg rounded-lg",children:(0,b.jsx)("div",{className:"px-4 py-3 bg-transparent",children:(0,b.jsx)("div",{className:"flex flex-wrap items-center",children:(0,b.jsxs)("div",{className:"relative w-full max-w-full flex-grow flex-1",children:[(0,b.jsx)("h6",{className:"text-gray-500 mb-1 text-xs font-semibold",children:"NS \u5206\u91cf\u5f53\u524d\u503c"}),(0,b.jsx)("h2",{className:"text-gray-700 text-xl font-semibold",children:this.state.analysis.north_south})]})})})}),(0,b.jsx)("div",{className:"relative flex flex-col min-w-0 break-words bg-sky-100 w-full mb-4 shadow-lg rounded-lg",children:(0,b.jsx)("div",{className:"px-4 py-3 bg-transparent",children:(0,b.jsx)("div",{className:"flex flex-wrap items-center",children:(0,b.jsxs)("div",{className:"relative w-full max-w-full flex-grow flex-1",children:[(0,b.jsx)("h6",{className:"text-gray-500 mb-1 text-xs font-semibold",children:"\u5408\u6210\u5206\u91cf\u5f53\u524d\u503c"}),(0,b.jsx)("h2",{className:"text-gray-700 text-xl font-semibold",children:this.state.analysis.synthesis})]})})})})]})})})})]})})]})]}),(0,b.jsx)(x.default,{}),(0,b.jsx)(h.default,{})]})}}]),s}(c.Component)},9526:function(e,t,s){s.r(t),s.d(t,{default:function(){return d}});var a=s(5671),l=s(3144),r=s(136),n=s(516),i=s(7313),o=s(6287),c=s(6417),d=function(e){(0,r.Z)(s,e);var t=(0,n.Z)(s);function s(e){var l;return(0,a.Z)(this,s),(l=t.call(this,e)).state={copyright:o.default.frontend.copyright},l}return(0,l.Z)(s,[{key:"render",value:function(){return(0,c.jsx)("footer",{className:"fixed bottom-0 w-full bg-gray-200 text-gray-500",children:(0,c.jsx)("div",{className:"container mx-auto flex flex-wrap flex-col sm:flex-row",children:(0,c.jsxs)("div",{className:"container mx-auto py-2 px-4 flex flex-wrap flex-col sm:flex-row",children:[(0,c.jsx)("span",{className:"text-xs text-center mt-1 ml-4 md:ml-12 lg:ml-16 md:text-left",children:this.props.extra||"Constructing Real-time Seismic Network Ambitiously."}),(0,c.jsx)("span",{className:"text-sm inline-flex sm:ml-auto sm:mt-0 mt-2 justify-center sm:justify-start",children:this.state.copyright})]})})})}}]),s}(i.Component)},6887:function(e,t,s){s.r(t),s.d(t,{default:function(){return d}});var a=s(5671),l=s(3144),r=s(136),n=s(516),i=s(7313),o=s(644),c=s(6417),d=function(e){(0,r.Z)(s,e);var t=(0,n.Z)(s);function s(e){var l;return(0,a.Z)(this,s),(l=t.call(this,e)).state={navigation:l.props.navigation||""},l}return(0,l.Z)(s,[{key:"render",value:function(){return(0,c.jsx)("nav",{className:"flex px-5 py-3 text-gray-700  rounded-lg bg-gray-50 mb-6",children:(0,c.jsxs)("ol",{className:"inline-flex items-center space-x-1 md:space-x-3",children:[(0,c.jsx)("li",{className:"inline-flex items-center",children:(0,c.jsxs)("div",{onClick:function(){return(0,o.default)({dest:"/",replace:!1})},className:"cursor-pointer inline-flex items-center text-sm font-medium text-gray-700 hover:text-gray-900",children:[(0,c.jsx)("svg",{className:"w-4 h-4 mr-2",fill:"currentColor",viewBox:"0 0 20 20",xmlns:"http://www.w3.org/2000/svg",children:(0,c.jsx)("path",{d:"M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"})}),"\u4e3b\u9875"]})}),this.state.navigation.length>0&&(0,c.jsx)("li",{children:(0,c.jsxs)("div",{className:"flex items-center",children:[(0,c.jsx)("svg",{className:"w-6 h-6 text-gray-400",fill:"currentColor",viewBox:"0 0 20 20",xmlns:"http://www.w3.org/2000/svg",children:(0,c.jsx)("path",{fillRule:"evenodd",d:"M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z",clipRule:"evenodd"})}),(0,c.jsx)("div",{className:"ml-1 text-sm font-medium text-gray-700 hover:text-gray-900 md:ml-2",children:this.state.navigation})]})})]})})}}]),s}(i.Component)},3874:function(e,t,s){s.r(t),s.d(t,{default:function(){return c}});var a=s(5671),l=s(3144),r=s(136),n=s(516),i=s(7313),o=s(6417),c=function(e){(0,r.Z)(s,e);var t=(0,n.Z)(s);function s(e){var l;return(0,a.Z)(this,s),(l=t.call(this,e)).state={},l}return(0,l.Z)(s,[{key:"render",value:function(){return(0,o.jsx)("div",{className:this.props.className,children:(0,o.jsxs)("div",{className:"flex flex-col gap-y-2 font-bold",children:[(0,o.jsxs)("div",{className:"flex flex-row gap-2 font-bold text-lg",children:[this.props.icon,(0,o.jsx)("span",{children:this.props.title})]}),(0,o.jsx)("span",{className:"pl-3 text-md font-medium",children:this.props.text})]})})}}]),s}(i.Component)},6013:function(e,t,s){s.r(t),s.d(t,{default:function(){return d}});var a=s(5671),l=s(3144),r=s(136),n=s(516),i=s(7313),o=s(1892),c=s(6417),d=function(e){(0,r.Z)(s,e);var t=(0,n.Z)(s);function s(e){var l;return(0,a.Z)(this,s),(l=t.call(this,e)).state={scrollTop:!1},l}return(0,l.Z)(s,[{key:"componentDidMount",value:function(){var e=this;(0,o.registerEvents)({eventArray:[{trigger:"scroll",id:"scroller_scrollTop"}],onEventCallback:function(){window.scrollY>100?e.setState({scrollTop:!0}):e.setState({scrollTop:!1})}})}},{key:"componentWillUnmount",value:function(){(0,o.removeEvents)([{trigger:"scroll",id:"scroller_scrollTop"}])}},{key:"render",value:function(){return(0,c.jsx)("button",{onClick:function(){return window.scrollTo({top:0,behavior:"smooth"})},className:"".concat(this.state.scrollTop?"inline-block":"hidden"," fixed p-3 bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded-full shadow-md hover:bg-purple-700 hover:shadow-lg focus:purple-red-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out bottom-24 right-5"),children:(0,c.jsx)("svg",{"aria-hidden":"true",focusable:"false",className:"w-4 h-4",role:"img",xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 448 512",children:(0,c.jsx)("path",{fill:"currentColor",d:"M34.9 289.5l-22.2-22.2c-9.4-9.4-9.4-24.6 0-33.9L207 39c9.4-9.4 24.6-9.4 33.9 0l194.3 194.3c9.4 9.4 9.4 24.6 0 33.9L413 289.4c-9.5 9.5-25 9.3-34.3-.4L264 168.6V456c0 13.3-10.7 24-24 24h-32c-13.3 0-24-10.7-24-24V168.6L69.2 289.1c-9.3 9.8-24.8 10-34.3.4z"})})})}}]),s}(i.Component)},8395:function(e,t,s){s.r(t),s.d(t,{default:function(){return h}});var a=s(5671),l=s(3144),r=s(136),n=s(516),i=s(7313),o=s(6287),c=s(2135),d=s(6417),h=function(e){(0,r.Z)(s,e);var t=(0,n.Z)(s);function s(e){var l;return(0,a.Z)(this,s),(l=t.call(this,e)).state={isSidebarOpen:!1,sidebarList:o.default.sidebar,sidebarVersion:o.default.frontend.version,sidebarTitle:o.default.frontend.title,sidebarMark:l.props.sidebarMark},l}return(0,l.Z)(s,[{key:"componentDidMount",value:function(){var e=this;this.state.sidebarList.forEach((function(t){e.state.sidebarMark!==t.tag||(document.title="".concat(t.title," | ").concat(e.state.sidebarTitle))}))}},{key:"render",value:function(){var e=this;return(0,d.jsxs)(d.Fragment,{children:[(0,d.jsxs)("div",{className:"fixed w-full z-30 flex bg-gray-200 p-2 items-center justify-center h-16 px-10",children:[(0,d.jsx)("div",{className:"".concat(this.state.isSidebarOpen||"ml-10"," text-gray-800 transform ease-in-out duration-500 flex-none h-full flex items-center justify-center text-lg font-bold"),children:this.state.sidebarTitle}),(0,d.jsx)("div",{className:"grow h-full flex items-center justify-center "}),(0,d.jsx)("div",{className:"flex-none h-full text-center flex items-center justify-center text-gray-500",children:(0,d.jsxs)("div",{className:"flex space-x-1 items-center lg:px-10",children:[(0,d.jsx)("div",{className:"flex-none flex justify-center",children:(0,d.jsx)("svg",{xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 512 512",className:"w-4 h-4",fill:"currentColor",children:(0,d.jsx)("path",{d:"M320 488c0 9.5-5.6 18.1-14.2 21.9s-18.8 2.3-25.8-4.1l-80-72c-5.1-4.6-7.9-11-7.9-17.8s2.9-13.3 7.9-17.8l80-72c7-6.3 17.2-7.9 25.8-4.1s14.2 12.4 14.2 21.9v40h16c35.3 0 64-28.7 64-64V153.3C371.7 141 352 112.8 352 80c0-44.2 35.8-80 80-80s80 35.8 80 80c0 32.8-19.7 61-48 73.3V320c0 70.7-57.3 128-128 128H320v40zM456 80a24 24 0 1 0 -48 0 24 24 0 1 0 48 0zM192 24c0-9.5 5.6-18.1 14.2-21.9s18.8-2.3 25.8 4.1l80 72c5.1 4.6 7.9 11 7.9 17.8s-2.9 13.3-7.9 17.8l-80 72c-7 6.3-17.2 7.9-25.8 4.1s-14.2-12.4-14.2-21.9V128H176c-35.3 0-64 28.7-64 64V358.7c28.3 12.3 48 40.5 48 73.3c0 44.2-35.8 80-80 80s-80-35.8-80-80c0-32.8 19.7-61 48-73.3V192c0-70.7 57.3-128 128-128h16V24zM56 432a24 24 0 1 0 48 0 24 24 0 1 0 -48 0z"})})}),(0,d.jsx)("div",{className:"md:block text-sm md:text-md",children:this.state.sidebarVersion})]})})]}),(0,d.jsxs)("aside",{className:"".concat(this.state.isSidebarOpen?"translate-x-none":"-translate-x-48"," w-60 fixed transition transform ease-in-out duration-1000 z-50 flex h-screen bg-gray-800"),children:[(0,d.jsx)("div",{className:"".concat(this.state.isSidebarOpen?"translate-x-0":"translate-x-24 scale-x-0"," w-full -right-6 transition transform ease-in duration-300 flex items-center justify-between border-4 border-white absolute top-2 rounded-full h-12"),children:(0,d.jsx)("div",{className:"flex items-center space-x-3 group bg-gradient-to-r from-indigo-500 via-purple-500 to-purple-500 pl-16 pr-6 py-2 rounded-full text-white",children:(0,d.jsx)("div",{className:"transform ease-in-out duration-300 mr-16 font-bold",children:"\u9762\u677f\u83dc\u5355"})})}),(0,d.jsx)("div",{onClick:function(){return e.setState({isSidebarOpen:!e.state.isSidebarOpen})},className:"-right-6 cursor-pointer transition transform ease-in-out duration-500 flex border-4 border-white bg-[#1E293B] hover:bg-purple-500 absolute top-2 p-3 rounded-full text-white hover:rotate-45",children:(0,d.jsx)("svg",{xmlns:"http://www.w3.org/2000/svg",fill:"none",viewBox:"0 0 24 24",strokeWidth:3,stroke:"currentColor",className:"w-4 h-4",children:(0,d.jsx)("path",{strokeLinecap:"round",strokeLinejoin:"round",d:"M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z"})})}),(0,d.jsx)("div",{className:"".concat(this.state.isSidebarOpen?"flex":"hidden"," text-white mt-20 flex-col space-y-2 w-full h-[calc(100vh)]"),children:this.state.sidebarList.map((function(t,s){return(0,d.jsxs)(c.rU,{to:t.link,className:"".concat(e.state.sidebarMark===t.tag?"text-purple-500":"text-white hover:text-purple-500"," cursor-pointer hover:ml-4 w-full bg-[#1E293B] p-2 pl-8 rounded-full transform ease-in-out duration-300 flex flex-row items-center space-x-3"),children:[t.icon,(0,d.jsx)("div",{children:t.title})]},s)}))}),(0,d.jsx)("div",{className:"".concat(this.state.isSidebarOpen?"hidden":"flex"," mt-20 flex-col space-y-2 w-full h-[calc(100vh)]"),children:this.state.sidebarList.map((function(t,s){return(0,d.jsx)(c.rU,{to:t.link,className:"".concat(e.state.sidebarMark===t.tag?"text-purple-500":"text-white hover:text-purple-500"," cursor-pointer justify-end pr-5 w-full bg-[#1E293B] p-3 rounded-full transform ease-in-out duration-300 flex"),children:t.icon},s)}))})]})]})}}]),s}(i.Component)},3481:function(e,t,s){s.r(t);t.default=function(e){var t=e.url,s=e.type,a=e.onOpenCallback,l=e.onMessageCallback,r=e.onCloseCallback,n=e.onErrorCallback,i=new WebSocket(t);return i.onopen=a,i.onmessage=l,i.onclose=r,i.onerror=n,i.binaryType=s,i}},9635:function(e,t,s){s.r(t);t.default=function(e){var t=e.host,s=e.port,a=e.api,l=e.version,r=e.tls,n=e.type,i="".concat(t,":").concat(s,"/api/").concat(l,"/").concat(a);switch(n){case"http":return r?"https://".concat(i):"http://".concat(i);case"websocket":return r?"wss://".concat(i):"ws://".concat(i);default:return null}}},2529:function(e,t,s){s.r(t);t.default=function(e){var t=e.getFullYear(),s=(e.getMonth()+1).toString().padStart(2,"0"),a=e.getDate().toString().padStart(2,"0"),l=e.getHours().toString().padStart(2,"0"),r=e.getMinutes().toString().padStart(2,"0"),n=e.getSeconds().toString().padStart(2,"0");return"".concat(t,"-").concat(s,"-").concat(a," ").concat(l,":").concat(r,":").concat(n)}}}]);