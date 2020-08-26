function main() {
  console.log("main function");
  console.log("ajax request to the resource which will require cors enabled");
  const username = new Date().toISOString();
  const mnemonic = "scissors system judge reveal slogan rice option body bronze insane evolve matter";

  $.ajax({
    type: 'POST',
    dataType: "json",
    url: "v0/user/signup?user=" + username + "&password=test&mnemonic=" + mnemonic,
    success: function (data) {
      console.log("log response on success");
      console.log(data);
    }
  });
}