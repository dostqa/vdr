from enum import Enum


class TypesPers(str, Enum):
    passport = "PASSPORT"
    inn = "INN"
    snils = "SNILS"
    phone = "PHONE"
    email = "EMAIL"
    address = "ADDRESS"
