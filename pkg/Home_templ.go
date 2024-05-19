// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.680
package goob

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func HomeTemplate() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<html class=\"h-full\"><head><title>Goob</title><link href=\"public/index.css\" rel=\"stylesheet\"></head><body class=\"h-full flex flex-col \">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = TopBarComponent().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"flex justify-around items-center flex-grow\"><a href=\"/newroom\" class=\"border-2 flex p-5 justify-center items-center rounded-md max-h-[33%] min-h-fit h-[300px] w-[300px] border-black shadow-lg hover:bg-black hover:text-white text-3xl hover:shadow-xl hover:scale-105 transition-all \"><button>Create new room</button></a><form action=\"/room\" class=\"border-2 flex flex-col p-5 gap-3 justify-around items-center rounded-md max-h-[33%] min-h-fit h-[300px] w-[300px] border-black shadow-lg  mb-0\"><p class=\"flex-grow flex items-center text-3xl\">Join room</p><input required type=\"text\" name=\"id\" class=\"border-2 rounded-sm w-5/6 h-6 border-gray-400 text-center\"> <input type=\"submit\" value=\"Go\" class=\"border-2 hover:bg-black hover:text-white flex-grow rounded-md border-black w-5/6 shadow-lg hover:shadow-xl hover:scale-105 transition-all \"></form></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
