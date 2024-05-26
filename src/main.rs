use axum::{http::StatusCode, routing::post, Extension, Json, Router};
use google_cloud_storage::{
    client::{Client, ClientConfig},
    sign::SignedURLOptions,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tower::ServiceBuilder;

struct AppState {
    storage_client: Client,
}

#[tokio::main]
async fn main() {
    // Create and authenticate storage client, then set it as shared state
    let config = ClientConfig::default().with_auth().await.unwrap();
    let storage_client = Client::new(config);
    let shared_state = Arc::new(AppState { storage_client });

    // Routes
    let app = Router::new()
        .route("/", post(handle_signed_url))
        .layer(ServiceBuilder::new().layer(Extension(shared_state)));

    // Run server
    let listener = tokio::net::TcpListener::bind("127.0.0.1:8080")
        .await
        .unwrap();
    println!("listening on {}", listener.local_addr().unwrap());
    axum::serve(listener, app.into_make_service())
        .await
        .unwrap();
}

#[derive(Deserialize)]
struct SignedUrlRequest {
    object: String,
}

#[derive(Serialize)]
struct SignUrlResponse {
    url: String,
}

async fn handle_signed_url(
    Extension(state): Extension<Arc<AppState>>,
    Json(payload): Json<SignedUrlRequest>,
) -> Result<Json<SignUrlResponse>, StatusCode> {
    let bucket = "symposia-dev-bucket";

    let object = payload.object;
    let storage_client = &state.storage_client;
    let url = storage_client
        .signed_url(bucket, &object, None, None, SignedURLOptions::default())
        .await;
    if let Err(_) = url {
        return Err(StatusCode::INTERNAL_SERVER_ERROR);
    }
    let response = SignUrlResponse { url: url.unwrap() };
    Ok(Json(response))
}
