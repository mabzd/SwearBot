lock '3.4.0'

application_files = [
  'bin/log.txt',
  'bin/mods/settings.json',
  'bin/mods/modswears/stats.json',
  'bin/mods/modswears/swears.txt']

downloadable_files = [
  'bin/token.txt',
  'bin/mods/config.json',
  'bin/mods/modchoice/config.json',
  'bin/mods/modmention/config.json',
  'bin/mods/modicm/config.json',
  'bin/mods/modswears/config.json']

linked_files = application_files + downloadable_files

set :application, "swearbot"
set :repo_url, "git@github.com:mabzd/SwearBot.git"

set :keep_releases, 5
set :deploy_to, "/var/go/swearbot"
set :linked_files, linked_files

namespace :app do
  task :compile do
    on roles(:app) do
      execute "cd #{current_path} && chmod +rwx ./compile.sh"
      execute "cd #{current_path} && GOPATH=$HOME/go ./compile.sh"
    end
  end

  %i(start stop).each do |command|
    task command do
      on roles(:app) do
        execute "/etc/init.d/swbot #{command}"
      end
    end
  end

  task :version do
    on roles(:app) do
      execute "chmod +x #{current_path}/version.sh"
      execute "cd #{repo_path} && #{current_path}/version.sh > #{current_path}/bin/version.txt"
    end
  end

  task :download_shared do
    on roles(:app) do
      FileUtils::mkdir_p "./shared"
      downloadable_files.each do |file_path|
        dir = File.dirname(file_path)
        FileUtils::mkdir_p "./shared/#{dir}"
        download! "#{shared_path}/#{file_path}", "./shared/#{dir}"
      end
    end
  end

  task :upload_shared do
    on roles(:app) do
      if File.directory?("./shared/bin")
        upload! "./shared/bin", "#{shared_path}", :recursive => true
      end
    end
  end

  task :clean_shared do
    on roles(:app) do
      FileUtils::rm_rf "./shared"
    end
  end
end

after "deploy", "app:compile"
after "deploy", "app:version"
after "deploy", "app:upload_shared"
after "deploy", "app:stop"
after "deploy", "app:start"
after "deploy", "app:clean_shared"
