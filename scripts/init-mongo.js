// scripts/init-mongo.js
db = db.getSiblingDB('nutrient_db');

// Foods collection indexes
db.foods.createIndex({ "searchTerms": "text" });
db.foods.createIndex({ "createdBy": 1, "visibility": 1 });
db.foods.createIndex({ "category": 1 });
db.foods.createIndex({ "source": 1 });

// Meal plans collection indexes
db.meal_plans.createIndex({ "userId": 1, "startDate": -1 });
db.meal_plans.createIndex({ "userId": 1, "planType": 1 });
db.meal_plans.createIndex({ "userId": 1, "status": 1 });

// Meal templates collection indexes
db.meal_templates.createIndex({ "userId": 1, "mealType": 1 });
db.meal_templates.createIndex({ "userId": 1, "isPublic": 1 });

// Shopping lists collection indexes
db.shopping_lists.createIndex({ "userId": 1, "mealPlanId": 1 });
db.shopping_lists.createIndex({ "userId": 1, "status": 1 });

// Users collection indexes
db.users.createIndex({ "email": 1 }, { unique: true });

print('Indexes created successfully');
