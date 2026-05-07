from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from typing import List
from datetime import datetime, timedelta
from app.core.database import get_db
from app.models import database as db_models
from app.models import schemas

router = APIRouter(prefix="/api/v1", tags=["events"])

@router.get("/events", response_model=List[schemas.EventResponse])
async def get_events(
    skip: int = 0,
    limit: int = 100,
    severity: str = None,
    event_type: str = None,
    db: Session = Depends(get_db)
):
    """Get security events with optional filters"""
    query = db.query(db_models.Event)
    
    if severity:
        query = query.filter(db_models.Event.severity == severity)
    
    if event_type:
        query = query.filter(db_models.Event.event_type == event_type)
    
    events = query.order_by(
        db_models.Event.created_at.desc()
    ).offset(skip).limit(limit).all()
    
    return [schemas.EventResponse.from_orm(e) for e in events]

@router.post("/events", response_model=schemas.EventResponse)
async def create_event(event: schemas.EventCreate, db: Session = Depends(get_db)):
    """Create a new security event"""
    db_event = db_models.Event(**event.model_dump())
    db.add(db_event)
    db.commit()
    db.refresh(db_event)
    return db_event

@router.get("/events/stats")
async def get_event_stats(db: Session = Depends(get_db)):
    """Get event statistics"""
    yesterday = datetime.utcnow() - timedelta(days=1)
    
    total_events = db.query(db_models.Event).count()
    events_24h = db.query(db_models.Event).filter(
        db_models.Event.created_at >= yesterday
    ).count()
    
    critical_events = db.query(db_models.Event).filter(
        db_models.Event.severity == "critical"
    ).count()
    
    warning_events = db.query(db_models.Event).filter(
        db_models.Event.severity == "warning"
    ).count()
    
    return {
        "total_events": total_events,
        "events_24h": events_24h,
        "critical_events": critical_events,
        "warning_events": warning_events
    }
