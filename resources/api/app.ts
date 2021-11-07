import {
    APIGatewayProxyEventV2,
    APIGatewayProxyResultV2,
    Context,
} from "https://deno.land/x/lambda/mod.ts";

import { opine } from "https://deno.land/x/opine@1.9.0/mod.ts";


export async function configureApp() {
    const app = opine();

    app.get("/", function (req, res) {
        res.send("Hello World");
    });
}


export async function handler(
    _event: APIGatewayProxyEventV2,
    _context: Context,
): Promise<APIGatewayProxyResultV2> {

    const app = await configureApp();

    return {
        statusCode: 200,
        headers: { "content-type": "text/html;charset=utf8" },
        body: `Welcome to deno ${Deno.version.deno} 🦕`,
    };
}