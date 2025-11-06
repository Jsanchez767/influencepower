# Deployment Checklist ✅

## Database (Supabase) - COMPLETED ✅

- [x] Schema applied (`schema.sql`)
- [x] Metrics schema applied (`schema_metrics.sql`)
- [x] Officials populated (53 officials)
- [x] City API data synced (107 bodies, 227 persons, 100 matters)
- [x] Metrics calculated for all officials

## Backend (Render) - ACTION REQUIRED ⚠️

### Deployment Status
Your backend should be deployed at: **`influencepower-backend.onrender.com`** (or similar)

### Environment Variables to Set in Render Dashboard

1. Go to your Render dashboard: https://dashboard.render.com
2. Select your `influencepower-backend` service
3. Go to "Environment" tab
4. Add these environment variables:

```
SUPABASE_URL=<your-supabase-project-url>
SUPABASE_ANON_KEY=<your-supabase-anon-key>
SUPABASE_SERVICE_ROLE_KEY=<your-supabase-service-role-key>
PORT=8080
```

**Where to find Supabase credentials:**
- Go to your Supabase project: https://supabase.com/dashboard/project/YOUR_PROJECT
- Settings → API
- Copy the values:
  - `Project URL` → SUPABASE_URL
  - `anon public` key → SUPABASE_ANON_KEY  
  - `service_role` key → SUPABASE_SERVICE_ROLE_KEY (⚠️ Keep secret!)

### After Setting Environment Variables
1. Render will automatically redeploy
2. Wait 2-3 minutes for build to complete
3. Test the health endpoint: `https://your-app.onrender.com/api/v1/health`

### Expected Response
```json
{
  "status": "healthy",
  "message": "InfluencePower API is running"
}
```

## Frontend (Vercel) - COMPLETED ✅

- [x] Code pushed to GitHub
- [x] Vercel auto-deploys from main branch
- [x] Environment variable `VITE_API_URL` set to backend URL
- [x] Build succeeds
- [x] Deployed to production

### Vercel Environment Variable
If not already set:
1. Go to Vercel dashboard: https://vercel.com/dashboard
2. Select your project
3. Settings → Environment Variables
4. Add: `VITE_API_URL` = `https://influencepower-backend.onrender.com/api/v1`
5. Redeploy if needed

## Testing After Deployment

### 1. Test Backend Health
```bash
curl https://your-backend-url.onrender.com/api/v1/health
```
Expected: `{"status":"healthy","message":"InfluencePower API is running"}`

### 2. Test Officials Endpoint
```bash
curl https://your-backend-url.onrender.com/api/v1/officials | jq '.[0]'
```
Expected: JSON with official data

### 3. Test Metrics Endpoint
```bash
curl https://your-backend-url.onrender.com/api/v1/officials/1/metrics | jq
```
Expected: JSON with metrics data

### 4. Test Frontend
1. Visit your Vercel URL (e.g., `https://influencepower.vercel.app`)
2. Click on an official card
3. Profile page should load with:
   - Official photo (or placeholder)
   - Real metrics data
   - Legislative impact scores
   - Voting records
   - Ward comparisons

## Troubleshooting

### Backend Not Responding
- Check Render logs: Dashboard → Your Service → Logs
- Verify environment variables are set
- Check build succeeded
- Ensure PORT is set to 8080

### "Failed to load official" Error
- Check browser console (F12)
- Verify VITE_API_URL in Vercel points to correct backend
- Check CORS is enabled (it should be)
- Verify backend is responding to API calls

### Metrics Not Loading
- Verify metrics schema was applied in Supabase
- Run calculate_metrics.go script if not already run
- Check backend logs for database connection errors

### Images Not Loading
- Default placeholder images will show
- To use real headshots: run `update_headshots.go` script
- Images will fallback to placeholders if unavailable

## Next Steps (Optional Improvements)

### Automatic Metrics Updates
Set up a cron job on Render:
1. Create a new "Cron Job" service in Render
2. Set schedule: `0 2 * * *` (runs daily at 2 AM)
3. Command: `cd backend && go run scripts/calculate_metrics.go`

### Headshots Sync
Run once manually:
```bash
cd backend
go run scripts/update_headshots.go
```

### Monitor Performance
- Set up Render metrics monitoring
- Enable Supabase connection pooling if needed
- Add caching layer (Redis) for frequently accessed data

## Quick Reference

### Important URLs
- **Frontend**: https://influencepower.vercel.app (or your Vercel domain)
- **Backend**: https://influencepower-backend.onrender.com (or your Render domain)
- **Supabase**: https://supabase.com/dashboard/project/YOUR_PROJECT
- **GitHub**: https://github.com/Jsanchez767/influencepower

### Key Files
- Backend API: `backend/handlers/handlers.go`
- Routes: `backend/api/routes.go`
- Frontend Profile: `frontend/src/components/ProfilePage.tsx`
- Database Schema: `backend/db/schema_metrics.sql`

## Status Summary

✅ = Complete
⚠️ = Action Required
⏳ = Optional/Future

- ✅ Database schema applied
- ✅ Metrics tables created
- ✅ Officials data populated
- ✅ Metrics calculated
- ⚠️ **Backend environment variables need to be set in Render**
- ✅ Frontend deployed to Vercel
- ✅ Code pushed to GitHub
- ⏳ Headshots sync (optional)
- ⏳ Automated metrics updates (optional)

---

**Most Important Next Step**: Set the Supabase environment variables in your Render dashboard! 🚀
