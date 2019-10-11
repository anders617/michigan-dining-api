CREATE TABLE DiningHalls (
  DiningHallId STRING(37) NOT NULL,
  Name STRING(128) NOT NULL,
  Campus STRING(128) NOT NULL,
  BuildingName STRING(128),
  AddressCity STRING(128),
  AddressPostalCode STRING(128),
  AddressState STRING(128),
  AddressStreet1 STRING(128),
  AddressStreet2 STRING(128),
  Type STRING(128),
  SortPosition INT64,
) PRIMARY KEY(DiningHallId);

CREATE TABLE DayHours (
  DiningHallId STRING(37) NOT NULL,
  Date DATE NOT NULL,
  Hour STRING(1024) NOT NULL,
) PRIMARY KEY (DiningHallId, Date),
INTERLEAVE IN PARENT DiningHalls ON DELETE CASCADE;

CREATE TABLE DayEvents (
  DiningHallId STRING(37) NOT NULL,
  Date DATE NOT NULL,
  DayStart STRING(128),
  TimeStart TIMESTAMP NOT NULL,
  DayEnd STRING(128),
  TimeEnd TIMESTAMP,
  Title STRING(128),
  Description STRING(128),
  MapLink STRING(128),
) PRIMARY KEY (DiningHallId, Date, TimeStart),
INTERLEAVE IN PARENT DiningHalls ON DELETE CASCADE;
