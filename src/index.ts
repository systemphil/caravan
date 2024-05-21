import express, { type Express, type Request, type Response } from "express";
import dotenv from "dotenv";
import { gcGenerateReadSignedUrl } from "./bucket";

dotenv.config();

const app: Express = express();
const port = process.env.PORT || 3000;

app.get("/", async (req: Request, res: Response) => {
    const fileName = "VID_20200103_135115.mp4";
    const id = "cluvqhyly0007uwfdmg2hn33a";
    const url = await gcGenerateReadSignedUrl({ fileName, id });
    res.send(url);
});

app.get("/:id/:filename", async (req: Request, res: Response) => {
    const fileName = req.params.filename;
    const id = req.params.id;
    const url = await gcGenerateReadSignedUrl({ fileName, id });
    res.send(url);
});

app.listen(port, () => {
    console.log(`Server listening on port ${port}`);
});
