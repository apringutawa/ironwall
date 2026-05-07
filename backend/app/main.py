from fastapi import FastAPI, Depends
from fastapi.middleware.cors import CORSMiddleware
from app.core.config import settings
from app.core.database import engine, Base
from app.api import status, events, scan, firewall, alerts

Base.metadata.create_all(bind=engine)

app = FastAPI(
    title=settings.APP_NAME,
    version=settings.VERSION,
    docs_url="/docs",
    redoc_url="/redoc"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(status.router)
app.include_router(events.router)
app.include_router(scan.router)
app.include_router(firewall.router)
app.include_router(alerts.router)

@app.get("/")
async def root():
    return {"name": settings.APP_NAME, "version": settings.VERSION}

@app.get("/health")
async def health():
    return {"status": "healthy"}
