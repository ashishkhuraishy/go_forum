ALTER TABLE "likes" ADD CONSTRAINT "post_user_unique" UNIQUE ("user_id", "post_id");


