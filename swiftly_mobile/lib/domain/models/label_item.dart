import 'dart:ui';

import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:uuid/uuid.dart';

part 'label_item.freezed.dart';

@freezed
class LabelItem with _$LabelItem{
  const factory LabelItem({
    required String id,
    required String cardId,
    required String userId,
    required String title,
    required Color color,
  }) = _LabelItem;

  factory LabelItem.create({
    required String cardId,
    required String userId,
    required String title,
    required Color color,
  }) => LabelItem(
    id: const Uuid().v4(),
    cardId: cardId,
    userId: userId,
    title: title,
    color: color,
  );
}