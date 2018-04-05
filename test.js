const puppeteer = require('puppeteer');

   

const tasks = Array()
var page = null

const errHander = function (reason) {
    console.log(reason)
}

!function(){
    puppeteer.launch().then((browser) => {
        browser.newPage().then(
            p => page = p
        ).catch(errHander)
    }).catch(errHander)  
}()

async function getContent(url, response) {

    // const browser = await puppeteer.launch();
    // const page = await browser.newPage();

    // res = await page.goto(url)
    // if (!res.ok()) {
    //     return "";
    // }
    // return await res.text();
    await page.goto(url)
    return page.url()
}


setTimeout(async ()=>{
    console.log(await getContent("http://music.163.com"))
    console.log(await getContent("http://qq.com"))
    console.log(await getContent("http://baidu.com"))
}, 3000)
