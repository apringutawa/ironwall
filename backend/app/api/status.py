from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from typing import List
from datetime import datetime, timedelta
from app.core.database import get_db
from app.models import database as db_models
from app.models import schemas

router = APIRouter(prefix="/api/v1", tags=["status"])

@router.get("/status", response_model=schemas.SystemStatus)
async def get_system_status(db: Session = Depends(get_db)):
    """Get overall system status and health"""
    
    # Get all modules
    modules = db.query(db_models.Module).all()
    
    # Count blocked IPs (mock data for now)
    blocked_ips_count = 147
    
    # Count events in last 24h
    yesterday = datetime.utcnow() - timedelta(days=1)
    events_24h = db.query(db_models.Event).filter(
        db_models.Event.created_at >= yesterday
    ).count()
    
    return schemas.SystemStatus(
        status="protected",
        uptime="5 days, 12 hours",
        cpu_usage=2.3,
        memory_usage=15.2,
        modules=[schemas.ModuleResponse.from_orm(m) for m in modules],
        blocked_ips_count=blocked_ips_count,
        events_24h=events_24h
    )

@router.get("/modules", response_model=List[schemas.ModuleResponse])
async def get_modules(db: Session = Depends(get_db)):
    """Get all security modules"""
    modules = db.query(db_models.Module).all()
    return [schemas.ModuleResponse.from_orm(m) for m in modules]

@router.post("/modules/{module_name}/enable")
async def enable_module(module_name: str, db: Session = Depends(get_db)):
    """Enable a security module"""
    module = db.query(db_models.Module).filter(
        db_models.Module.name == module_name
    ).first()
    
    if not module:
        raise HTTPException(status_code=404, detail="Module not found")
    
    module.enabled = True
    module.status = "active"
    db.commit()
    
    return {"message": f"Module {module_name} enabled"}

@router.post("/modules/{module_name}/disable")
async def disable_module(module_name: str, db: Session = Depends(get_db)):
    """Disable a security module"""
    module = db.query(db_models.Module).filter(
        db_models.Module.name == module_name
    ).first()
    
    if not module:
        raise HTTPException(status_code=404, detail="Module not found")
    
    module.enabled = False
    module.status = "inactive"
    db.commit()
    
    return {"message": f"Module {module_name} disabled"}
