class FileInfo {
  final String name;
  final String type; // image, folder, archive, file, pdf, doc
  final String? localPath;

  FileInfo(this.name, this.type, this.localPath);
  
  String get extension =>
      name.split('.').length > 1 ? name.split('.').last.toLowerCase() : '';
}