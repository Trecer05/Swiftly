import 'package:freezed_annotation/freezed_annotation.dart';

import '../../models/label_item.dart';

part 'user.freezed.dart';

@freezed
class User with _$User {
  const factory User({
    required String id,
    required String name,
    required String? lastName,
    required String? image,
    required LabelItem? role,
  }) = _User;

  factory User.create({
    required String id,
    String? name,
    String? lastName,
    String? image,
    LabelItem? role,
  }) => User(
    id: id,
    name: name ?? 'Новый пользователь',
    lastName: lastName,
    image: image,
    role: role,
  );
}
