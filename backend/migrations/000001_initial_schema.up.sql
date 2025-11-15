-- users and identities
CREATE TABLE IF NOT EXISTS users (
                                     id UUID PRIMARY KEY,
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_identities (
                                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                               user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                               provider TEXT NOT NULL,
                                               provider_user_id TEXT NOT NULL,
                                               data JSONB,
                                               UNIQUE(provider, provider_user_id)
);

CREATE TABLE IF NOT EXISTS user_profiles (
                                             user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                                             display_name TEXT,
                                             email TEXT,
                                             phone TEXT,
                                             tz TEXT DEFAULT 'UTC',
                                             locale TEXT DEFAULT 'ru',
                                             saved_fields JSONB DEFAULT '{}'::jsonb
);

-- events
CREATE TABLE IF NOT EXISTS events (
                                      id UUID PRIMARY KEY,
                                      owner_id UUID NOT NULL REFERENCES users(id),
                                      title TEXT NOT NULL,
                                      description TEXT,
                                      visibility TEXT NOT NULL DEFAULT 'public',
                                      status TEXT NOT NULL DEFAULT 'draft',
                                      starts_at TIMESTAMPTZ NOT NULL,
                                      ends_at TIMESTAMPTZ,
                                      tz TEXT NOT NULL DEFAULT 'UTC',
                                      location TEXT,
                                      online_url TEXT,
                                      capacity INT DEFAULT 0,
                                      waitlist_enabled BOOLEAN DEFAULT true,
                                      settings JSONB DEFAULT '{}'::jsonb,
                                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS event_series (
                                            id UUID PRIMARY KEY,
                                            event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                            rrule TEXT NOT NULL,
                                            exdates TIMESTAMPTZ[],
                                            until TIMESTAMPTZ,
                                            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                            updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS event_roles (
                                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                           event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                           user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                           role TEXT NOT NULL,
                                           UNIQUE(event_id, user_id)
);

-- forms
CREATE TABLE IF NOT EXISTS forms (
                                     id UUID PRIMARY KEY,
                                     event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                     version INT NOT NULL DEFAULT 1,
                                     schema JSONB NOT NULL,
                                     rules JSONB,
                                     active BOOLEAN DEFAULT true,
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS form_responses (
                                              id UUID PRIMARY KEY,
                                              form_id UUID NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
                                              user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                              status TEXT NOT NULL DEFAULT 'draft',
                                              answers JSONB NOT NULL DEFAULT '{}'::jsonb,
                                              created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                              updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                              UNIQUE(form_id, user_id)
);

-- registrations
CREATE TABLE IF NOT EXISTS registrations (
                                             id UUID PRIMARY KEY,
                                             event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                             user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                             ticket_type_id UUID,
                                             status TEXT NOT NULL DEFAULT 'going',
                                             source TEXT,
                                             utm JSONB,
                                             created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                             updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                             UNIQUE(event_id, user_id)
);

CREATE TABLE IF NOT EXISTS waitlist (
                                        id UUID PRIMARY KEY,
                                        event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                        UNIQUE(event_id, user_id)
);

-- campaigns and deliveries
CREATE TABLE IF NOT EXISTS campaigns (
                                         id UUID PRIMARY KEY,
                                         event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                         name TEXT NOT NULL,
                                         segment TEXT NOT NULL,
                                         content JSONB NOT NULL,
                                         schedule_at TIMESTAMPTZ,
                                         status TEXT NOT NULL DEFAULT 'pending',
                                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                         updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS deliveries (
                                          id UUID PRIMARY KEY,
                                          campaign_id UUID REFERENCES campaigns(id) ON DELETE SET NULL,
                                          channel TEXT NOT NULL,
                                          target_user_id UUID NOT NULL REFERENCES users(id),
                                          message_id TEXT,
                                          status TEXT NOT NULL DEFAULT 'pending',
                                          error TEXT,
                                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reminders (
                                         id UUID PRIMARY KEY,
                                         event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                         user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                         at TIMESTAMPTZ NOT NULL,
                                         type TEXT NOT NULL,
                                         status TEXT NOT NULL DEFAULT 'pending',
                                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- polls
CREATE TABLE IF NOT EXISTS polls (
                                     id UUID PRIMARY KEY,
                                     event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                     question TEXT NOT NULL,
                                     options JSONB NOT NULL,
                                     type TEXT NOT NULL DEFAULT 'single',
                                     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS poll_votes (
                                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                          poll_id UUID NOT NULL REFERENCES polls(id) ON DELETE CASCADE,
                                          user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                          option_key TEXT NOT NULL,
                                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                          UNIQUE(poll_id, user_id)
);

-- checkin
CREATE TABLE IF NOT EXISTS checkins (
                                        id UUID PRIMARY KEY,
                                        event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                        method TEXT NOT NULL,
                                        at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS qr_tokens (
                                         id UUID PRIMARY KEY,
                                         user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                         event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                         token_hash BYTEA NOT NULL UNIQUE,
                                         expires_at TIMESTAMPTZ NOT NULL,
                                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- analytics
CREATE TABLE IF NOT EXISTS event_views (
                                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                           event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
                                           user_id UUID REFERENCES users(id) ON DELETE SET NULL,
                                           source TEXT,
                                           utm JSONB,
                                           at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_log (
                                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                         actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
                                         entity TEXT NOT NULL,
                                         entity_id UUID NOT NULL,
                                         action TEXT NOT NULL,
                                         data JSONB,
                                         at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- subscriptions
CREATE TABLE IF NOT EXISTS subscriptions (
                                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                             user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                             type TEXT NOT NULL,
                                             enabled BOOLEAN DEFAULT true,
                                             meta JSONB
);

-- indexes
CREATE INDEX idx_events_owner ON events(owner_id);
CREATE INDEX idx_events_starts_at ON events(starts_at);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_visibility ON events(visibility);

CREATE INDEX idx_registrations_event ON registrations(event_id);
CREATE INDEX idx_registrations_user ON registrations(user_id);
CREATE INDEX idx_registrations_status ON registrations(status);

CREATE INDEX idx_forms_event ON forms(event_id);
CREATE INDEX idx_forms_active ON forms(active) WHERE active = true;

CREATE INDEX idx_form_responses_form ON form_responses(form_id);
CREATE INDEX idx_form_responses_user ON form_responses(user_id);

CREATE INDEX idx_polls_event ON polls(event_id);
CREATE INDEX idx_poll_votes_poll ON poll_votes(poll_id);

CREATE INDEX idx_checkins_event ON checkins(event_id);
CREATE INDEX idx_checkins_user ON checkins(user_id);

CREATE INDEX idx_qr_tokens_hash ON qr_tokens(token_hash);
CREATE INDEX idx_qr_tokens_expires ON qr_tokens(expires_at);

CREATE INDEX idx_campaigns_event ON campaigns(event_id);
CREATE INDEX idx_campaigns_schedule ON campaigns(schedule_at);
CREATE INDEX idx_campaigns_status ON campaigns(status);

CREATE INDEX idx_deliveries_campaign ON deliveries(campaign_id);
CREATE INDEX idx_deliveries_target_user ON deliveries(target_user_id);
CREATE INDEX idx_deliveries_status ON deliveries(status);

CREATE INDEX idx_reminders_event ON reminders(event_id);
CREATE INDEX idx_reminders_user ON reminders(user_id);
CREATE INDEX idx_reminders_at ON reminders(at);

CREATE INDEX idx_event_views_event ON event_views(event_id);
CREATE INDEX idx_event_views_at ON event_views(at);