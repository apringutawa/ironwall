from sqlalchemy import Column, Integer, String, Boolean, DateTime, JSON, ForeignKey
from sqlalchemy.orm import relationship
from datetime import datetime
from app.core.database import Base

class Server(Base):
    __tablename__ = "servers"
    
    id = Column(Integer, primary_key=True, index=True)
    hostname = Column(String, nullable=False)
    ip_address = Column(String, nullable=False)
    os_type = Column(String, nullable=False)
    os_version = Column(String, nullable=False)
    installed_at = Column(DateTime, default=datetime.utcnow)
    last_scan = Column(DateTime, nullable=True)
    
    events = relationship("Event", back_populates="server")

class Event(Base):
    __tablename__ = "events"
    
    id = Column(Integer, primary_key=True, index=True)
    server_id = Column(Integer, ForeignKey("servers.id"))
    event_type = Column(String, nullable=False)
    severity = Column(String, nullable=False)
    description = Column(String, nullable=False)
    details = Column(JSON, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)
    
    server = relationship("Server", back_populates="events")

class Module(Base):
    __tablename__ = "modules"
    
    id = Column(Integer, primary_key=True, index=True)
    name = Column(String, unique=True, nullable=False)
    enabled = Column(Boolean, default=True)
    status = Column(String, default="active")
    last_run = Column(DateTime, nullable=True)

class Config(Base):
    __tablename__ = "config"
    
    key = Column(String, primary_key=True)
    value = Column(String, nullable=False)

class Backup(Base):
    __tablename__ = "backups"
    
    id = Column(Integer, primary_key=True, index=True)
    file_path = Column(String, nullable=False)
    backup_path = Column(String, nullable=False)
    created_at = Column(DateTime, default=datetime.utcnow)
