class Article < ApplicationRecord
  validate :title_does_not_contain_balalaika_word

  private

  def title_does_not_contain_balalaika_word
    File.open(Rails.root.join("app/assets/javascripts/validators/articles.js")) do |file|
      cxt = V8::Context.new
      cxt.eval(file)

      unless cxt[:valid].call(attributes.symbolize_keys)
        errors.add(:title, "must not contain balalaika word")
      end
    end
  end
end
