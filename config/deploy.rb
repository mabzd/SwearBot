lock '3.4.0'

set :application, "swearbot"
set :repo_url, "git@github.com:mabzd/SwearBot.git"

set :keep_releases, 5
set :deploy_to, "/var/go/swearbot"
set :linked_files, [
  'bin/token.txt',
  'bin/log.txt',
  'bin/mods/config.json',
  'bin/mods/settings.json',
  'bin/mods/modchoice/config.json',
  'bin/mods/modmention/config.json',
  'bin/mods/modswears/config.json',
  'bin/mods/modswears/stats.json',
  'bin/mods/modswears/swears.txt']

namespace :app do
  task :compile do
    on roles(:app) do
      execute "cd #{current_path} && GOPATH=$HOME/go && ./compile.sh"
    end
  end

  %i(start stop).each do |command|
    task command do
      on roles(:app) do
        execute "/etc/init.d/swbot #{command}"
      end
    end
  end
end

after "deploy", "app:compile"
after "deploy", "app:stop"
after "deploy", "app:start"
