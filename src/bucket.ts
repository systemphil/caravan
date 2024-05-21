import { Storage, type GetSignedUrlConfig } from "@google-cloud/storage";

const storage = new Storage();

const PRIMARY_BUCKET_NAME = process.env.GCP_PRIMARY_BUCKET_NAME ?? "invalid";
const SECONDARY_BUCKET_NAME =
    process.env.GCP_SECONDARY_BUCKET_NAME ?? "invalid";

/**
 * Storage bucket reference instance. Call methods on this object to interact with the bucket.
 * @see https://cloud.google.com/nodejs/docs/reference/storage/latest
 */
export const primaryBucket = storage.bucket(PRIMARY_BUCKET_NAME);
export const secondaryBucket = storage.bucket(SECONDARY_BUCKET_NAME);

// export async function TestBucket() {
//     const sa = await storage.getProjectId();
//     console.log(sa);
//     const a = await storage.getServiceAccount();
//     console.log(a);
// }

export async function gcGenerateReadSignedUrl({
    fileName,
    id,
}: {
    fileName: string;
    id: string;
}) {
    const filePath = `video/${id}/${fileName}`;
    const options: GetSignedUrlConfig = {
        version: "v4",
        action: "read",
        expires: Date.now() + 15 * 60 * 1000, // 15 minutes
    };
    const [url] = await primaryBucket.file(filePath).getSignedUrl(options);
    return url;
}
